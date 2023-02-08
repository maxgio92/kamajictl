package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/cli-utils/pkg/kstatus/polling"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/api/konfig"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen/kustomization"
	runclient "github.com/fluxcd/pkg/runtime/client"
	"github.com/fluxcd/pkg/ssa"
)

// Apply is the equivalent of 'kubectl apply --server-side -f'.
// If the given manifest is a kustomization.yaml, then apply performs the equivalent of 'kubectl apply --server-side -k'.
func Apply(ctx context.Context, rcg genericclioptions.RESTClientGetter, opts *runclient.Options, root, manifestPath string) (string, error) {
	objs, err := readObjects(root, manifestPath)
	if err != nil {
		return "", err
	}

	if len(objs) == 0 {
		return "", fmt.Errorf("no Kubernetes objects found at: %s", manifestPath)
	}

	// Work around bugs in server-side apply upstream affecting Kubernetes < 1.22.
	if err := ssa.SetNativeKindsDefaults(objs); err != nil {
		return "", err
	}

	changeSet := ssa.NewChangeSet()

	// contains only CRDs and Namespaces
	var stageOne []*unstructured.Unstructured

	// contains all objects except for CRDs and Namespaces
	var stageTwo []*unstructured.Unstructured

	for _, u := range objs {
		if ssa.IsClusterDefinition(u) {
			stageOne = append(stageOne, u)
		} else {
			stageTwo = append(stageTwo, u)
		}
	}

	if len(stageOne) > 0 {
		cs, err := applySet(ctx, rcg, opts, stageOne)
		if err != nil {
			return "", err
		}
		changeSet.Append(cs.Entries)
	}

	if err := waitForSet(rcg, opts, changeSet); err != nil {
		return "", err
	}

	if len(stageTwo) > 0 {
		cs, err := applySet(ctx, rcg, opts, stageTwo)
		if err != nil {
			return "", err
		}
		changeSet.Append(cs.Entries)
	}

	return changeSet.String(), nil
}

func readObjects(root, manifestPath string) ([]*unstructured.Unstructured, error) {
	fi, err := os.Lstat(manifestPath)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() || !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("expected %q to be a file", manifestPath)
	}

	if isKustomizationFile(manifestPath) {
		resources, err := kustomization.BuildWithRoot(root, filepath.Dir(manifestPath))
		if err != nil {
			return nil, err
		}
		return ssa.ReadObjects(bytes.NewReader(resources))
	}

	ms, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer ms.Close()

	return ssa.ReadObjects(bufio.NewReader(ms))
}

func applySet(ctx context.Context, rcg genericclioptions.RESTClientGetter, opts *runclient.Options, objects []*unstructured.Unstructured) (*ssa.ChangeSet, error) {
	man, err := newManager(rcg, opts)
	if err != nil {
		return nil, err
	}

	return man.ApplyAll(ctx, objects, ssa.DefaultApplyOptions())
}

func waitForSet(rcg genericclioptions.RESTClientGetter, opts *runclient.Options, changeSet *ssa.ChangeSet) error {
	man, err := newManager(rcg, opts)
	if err != nil {
		return err
	}

	return man.WaitForSet(
		changeSet.ToObjMetadataSet(),
		ssa.WaitOptions{
			Interval: 2 * time.Second,
			Timeout:  time.Minute,
		},
	)
}

func newManager(rcg genericclioptions.RESTClientGetter, opts *runclient.Options) (*ssa.ResourceManager, error) {
	cfg, err := NewKubeConfigWithOptions(rcg)
	if err != nil {
		return nil, err
	}

	restMapper, err := rcg.ToRESTMapper()
	if err != nil {
		return nil, err
	}

	kubeClient, err := client.New(cfg, client.Options{Mapper: restMapper, Scheme: NewScheme()})
	if err != nil {
		return nil, err
	}

	kubePoller := polling.NewStatusPoller(kubeClient, restMapper, polling.Options{})

	return ssa.NewResourceManager(kubeClient, kubePoller, ssa.Owner{
		Field: "kamaji",
		Group: "kamaji.clastix.io",
	}), nil
}

func isKustomizationFile(path string) bool {
	base := filepath.Base(path)

	for _, v := range konfig.RecognizedKustomizationFileNames() {
		if base == v {
			return true
		}
	}

	return false
}
