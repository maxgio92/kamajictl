package install

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen"
	"github.com/fluxcd/flux2/v2/pkg/status"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/cli-utils/pkg/object"

	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/utils"
)

type Options struct {
	version       string
	manifestsPath string
	*options.CommonOptions
}

// NewInstallCommand returns the TCP command.
func NewInstallCommand(opts *options.CommonOptions) *cobra.Command {
	// Instantiate the command options.
	o := &Options{
		CommonOptions: opts,
	}

	// Instantiate the command.
	cmd := &cobra.Command{
		Use:           commandName,
		Short:         commandShortDescription,
		RunE:          o.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVar(&o.version, "version", defaultVersion, "the kamaji version to install")
	cmd.Flags().StringVar(&o.manifestsPath, "manifests-path", "manifests", "the path to the kamaji manifests")
	o.CommonOptions.AddFlags(cmd.Flags())

	return cmd
}

func (o *Options) Run(_ *cobra.Command, _ []string) error {
	// Validate CLI options.
	if err := o.CommonOptions.Validate(); err != nil {
		return errors.Wrap(err, "error validating tcp create command options")
	}

	// TODO: move this logic inside KamajiClient.
	tmpDir, err := manifestgen.MkdirTempAbs("", *o.KubeconfigOptions.Namespace)
	if err != nil {
		return errors.Wrap(err, "error generating manifests")
	}
	defer os.RemoveAll(tmpDir)

	if !isEmbeddedVersion(o.version) {
		return errors.New("the kamaji version is not supported")
	}

	if err := o.writeEmbeddedManifests(tmpDir); err != nil {
		return err
	}
	manifestsBase := tmpDir

	o.Logger.Actionf("applying resources")

	_, err = utils.Apply(
		o.Context,
		o.KubeconfigOptions,
		kubeclientOptions,
		tmpDir,
		filepath.Join(manifestsBase, kustomization),
	)
	if err != nil {
		return errors.Wrap(err, "install failed")
	}

	o.Logger.Successf("resources applied")

	kubeConfig, err := utils.NewKubeConfigWithOptions(o.KubeconfigOptions)
	if err != nil {
		return errors.Wrap(err, "install failed")
	}

	statusChecker, err := status.NewStatusChecker(kubeConfig, 5*time.Second, timeout, o.Logger)
	if err != nil {
		return errors.Wrap(err, "install failed")
	}

	componentRefs, err := buildComponentObjectRefs(*o.KubeconfigOptions.Namespace, components...)
	if err != nil {
		return errors.Wrap(err, "install failed")
	}

	o.Logger.Waitingf("verifying installation")

	if err := statusChecker.Assess(componentRefs...); err != nil {
		return errors.New("install failed")
	}

	o.Logger.Successf("Install finished")

	return nil
}

func isEmbeddedVersion(input string) bool {
	return input == defaultVersion
}

func (o *Options) writeEmbeddedManifests(dir string) error {
	manifests, err := fs.ReadDir(o.Embedded, "manifests")
	if err != nil {
		return err
	}
	for _, manifest := range manifests {
		data, err := fs.ReadFile(o.Embedded, path.Join("manifests", manifest.Name()))
		if err != nil {
			return fmt.Errorf("reading file failed: %w", err)
		}

		err = os.WriteFile(path.Join(dir, manifest.Name()), data, 0666)
		if err != nil {
			return fmt.Errorf("writing file failed: %w", err)
		}
	}
	return nil
}

func buildComponentObjectRefs(namespace string, components ...object.ObjMetadata) ([]object.ObjMetadata, error) {
	var objRefs []object.ObjMetadata

	for _, obj := range components {
		obj.Namespace = namespace
		objRefs = append(objRefs, obj)
	}

	return objRefs, nil
}
