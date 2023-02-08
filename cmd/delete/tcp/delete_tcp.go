package tcp

import (
	"github.com/maxgio92/kamajicli/internal/utils"
	"github.com/maxgio92/kamajicli/pkg/tcp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajicli/internal/options"
)

type Options struct {
	Client *tcp.TCPClient
	TCP    *tcp.TCPOptions
	*options.CommonOptions
}

// NewTCPCommand returns the TCP command.
func NewTCPCommand(opts *options.CommonOptions) *cobra.Command {
	// Instantiate the command options.
	o := &Options{
		&tcp.TCPClient{},
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

	o.TCP.AddCommonFlags(cmd.Flags())
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
		return errors.Wrap(err, "error validating tcp delete command options")
	}

	// Instantiate the Kubernetes client.
	kube, err := utils.NewKubeClient()
	if err != nil {
		return err
	}

	// Instantiate the TCP client.
	tcpClient, err := tcp.NewTCPClient(
		tcp.WithLogger(o.Logger),
		tcp.WithKubeClient(kube),
	)
	if err != nil {
		return err
	}

	o.Client = tcpClient

	// Set TCP options.
	o.TCP.Set(
		tcp.WithName(name),
	)

	// Delete the TCP.
	if err := o.Client.DeleteTCP(o.Context, o.TCP); err != nil {
		return errors.Wrap(err, "error deleting tenantcontrolplane")
	}

	o.Printer.Success.Printf("TenantControlPlane %s/%s deleted", o.TCP.Namespace, name)

	return nil
}
