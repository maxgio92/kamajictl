package tcp

import (
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/utils"
	"github.com/maxgio92/kamajictl/pkg/kamaji"
	"github.com/maxgio92/kamajictl/pkg/kamaji/tcp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Options struct {
	Client *kamaji.KamajiClient
	TCP    *tcp.TCPOptions
	*options.CommonOptions
}

// NewTCPCommand returns the TCP command.
func NewTCPCommand(opts *options.CommonOptions) *cobra.Command {
	// Instantiate the command options.
	o := &Options{
		&kamaji.KamajiClient{},
		tcp.NewTCPOptions(),
		opts,
	}

	// Instantiate the command.
	cmd := &cobra.Command{
		Use:           commandName,
		Short:         commandShortDescription,
		Args:          cobra.MinimumNArgs(1),
		RunE:          o.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	o.TCP.AddCreateFlags(cmd.Flags())
	o.CommonOptions.AddFlags(cmd.Flags())

	return cmd
}

func (o *Options) Run(_ *cobra.Command, args []string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// Validate CLI options.
	if err := o.CommonOptions.Validate(); err != nil {
		return errors.Wrap(err, "error validating tcp create command options")
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

	// Set TCP options.
	o.TCP.Set(
		tcp.WithName(name),
		tcp.WithNamespace(*o.KubeconfigOptions.Namespace),
	)

	// Create the TCP.
	if err := o.Client.CreateTCP(o.Context, o.TCP); err != nil {
		return errors.Wrap(err, "error creating tenantcontrolplane")
	}

	o.Printer.Success.Printf("TenantControlPlane %s/%s created\n", o.TCP.Namespace, name)

	return nil
}
