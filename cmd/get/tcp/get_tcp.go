package tcp

import (
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/output"
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
		RunE:          o.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

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
		return errors.Wrap(err, "error validating tcp get command options")
	}

	// Instantiate the Kubernetes client.
	kube, err := utils.NewKubeClient()
	if err != nil {
		return err
	}

	kc, err := kamaji.NewKamajiClient(
		kamaji.WithLogger(o.Logger),
		kamaji.WithKubeClient(kube),
	)
	if err != nil {
		return err
	}

	o.Client = kc

	// Set TCP options.
	o.TCP.Set(tcp.WithName(name))

	if !o.AllNamespaces {
		o.TCP.Set(tcp.WithNamespace(*o.KubeconfigOptions.Namespace))
	}

	// Get the TCPs.
	if name == "" {
		tcps, err := o.Client.ListTCPs(o.Context, o.TCP)
		if err != nil {
			return errors.Wrap(err, "error getting tenantcontrolplane")
		}

		if len(tcps.Items) > 0 {
			o.Printer.PrintTable(output.TCPListTable(tcps.Items...))

			return nil
		}
		o.Printer.Warningf("no tenantcontrolplane found")

		return nil
	}

	// Get the TCP.
	tcp, err := o.Client.GetTCP(o.Context, o.TCP)
	if err != nil {
		return errors.Wrap(err, "error getting tenantcontrolplane")
	}

	o.Printer.PrintTable(output.TCPListTable(*tcp))

	return nil
}
