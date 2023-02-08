package kubeconfig

import (
	"github.com/maxgio92/kamajicli/internal/options"
	"github.com/maxgio92/kamajicli/internal/utils"
	"github.com/maxgio92/kamajicli/pkg/tcp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Options struct {
	Client *tcp.TCPClient
	TCP    *tcp.TCPOptions
	*options.CommonOptions
}

// NewKubeconfigCommand returns the TCP command.
func NewKubeconfigCommand(opts *options.CommonOptions) *cobra.Command {

	// Instantiate the Kubernetes client.
	kube, err := utils.NewKubeClient()
	if err != nil {
		opts.Printer.ExitOnErr(err)
	}

	// Instantiate the TCP client.
	tcpClient, err := tcp.NewTCPClient(
		tcp.WithLogger(opts.Logger),
		tcp.WithKubeClient(kube),
	)
	if err != nil {
		opts.Printer.ExitOnErr(err)
	}

	// Instantiate the command options.
	o := &Options{
		tcpClient,
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
		SilenceUsage:  false,
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
		return errors.Wrap(err, "error validating tcp get command options")
	}

	// Set TCP options.
	o.TCP.Set(tcp.WithName(name))

	// Get the kubeconfig.
	kubeconfig, err := o.Client.GetKubeconfig(o.Context, o.TCP)
	if err != nil {
		return errors.Wrap(err, "error getting kubeconfig")
	}

	o.Printer.Raw(kubeconfig)

	return nil
}
