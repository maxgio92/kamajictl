package uninstall

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/utils"
	"github.com/maxgio92/kamajictl/pkg/kamaji"
)

type Options struct {
	Client *kamaji.KamajiClient
	*options.CommonOptions
}

// NewUninstallCommand returns the TCP command.
func NewUninstallCommand(opts *options.CommonOptions) *cobra.Command {
	// Instantiate the command options.
	o := &Options{
		&kamaji.KamajiClient{},
		opts,
	}

	// Instantiate the command.
	cmd := &cobra.Command{
		Use:           commandName,
		Short:         commandShortDescription,
		RunE:          o.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	o.CommonOptions.AddFlags(cmd.Flags())

	return cmd
}

func (o *Options) Run(_ *cobra.Command, args []string) error {
	// Validate CLI options.
	if err := o.CommonOptions.Validate(); err != nil {
		return errors.Wrap(err, "error validating tcp delete command options")
	}

	// Instantiate the Kubernetes client.
	kube, err := utils.NewKubeClient()
	if err != nil {
		return err
	}

	// Instantiate the Kamaji client.
	kc, err := kamaji.NewKamajiClient(
		kamaji.WithLogger(o.Logger),
		kamaji.WithKubeClient(kube),
	)
	if err != nil {
		return err
	}

	o.Client = kc

	//// Delete all the TCPs.
	if err := o.Client.PurgeTCPs(o.Context); err != nil {
		return errors.Wrap(err, "error cleaning up TCPs")
	}
	o.Logger.Successf("TenantControlPlanes deleted")

	o.Logger.Actionf("Deleting Kamaji controller")

	if err := kube.Delete(o.Context, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: *o.KubeconfigOptions.Namespace,
		}},
	); err != nil {
		return errors.Wrap(err, "error deleting Kamaji controller")
	}
	o.Logger.Successf("Kamaji controller uninstalled")

	o.Logger.Actionf("Deleting Kamaji CRDs")

	// Delete all the CRDs.
	if err := o.Client.PurgeCRDs(o.Context); err != nil {
		return errors.Wrap(err, "error deleting Kamaji CRDs")
	}
	o.Logger.Successf("Kamaji CRDs deleted")

	o.Logger.Actionf("Deleting Kamaji webhook configurations")

	// Delete the webhook configurations.
	if err := o.Client.PurgeWebhooks(o.Context); err != nil {
		return errors.Wrap(err, "error deleting Kamaji webhook configurations")
	}
	o.Logger.Successf("Kamaji webhook configurations deleted")

	o.Logger.Actionf("Deleting Kamaji RBAC resources")

	// Delete the webhook configurations.
	if err := o.Client.PurgeRBAC(o.Context); err != nil {
		return errors.Wrap(err, "error deleting Kamaji RBAC resoureces")
	}
	o.Logger.Successf("Kamaji RBAC resources deleted")

	o.Logger.Successf("Uninstall finished")

	return nil
}
