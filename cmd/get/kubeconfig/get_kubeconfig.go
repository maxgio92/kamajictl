package kubeconfig

import (
	"github.com/ghodss/yaml"
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

	// The format of the command's output.
	OutputFormat string
}

// NewKubeconfigCommand returns the TCP command.
func NewKubeconfigCommand(opts *options.CommonOptions) *cobra.Command {

	// Instantiate the command options.
	o := &Options{
		&kamaji.KamajiClient{},
		tcp.NewTCPOptions(),
		opts,
		"",
	}

	// Instantiate the command.
	cmd := &cobra.Command{
		Use:           commandName,
		Short:         commandShortDescription,
		RunE:          o.Run,
		SilenceErrors: true,
		SilenceUsage:  false,
	}

	o.CommonOptions.AddFlags(cmd.Flags())
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", `Output format. One of ("yaml", "table")`)

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
	o.TCP.Set(tcp.WithName(name))

	if !o.AllNamespaces {
		o.TCP.Set(tcp.WithNamespace(*o.KubeconfigOptions.Namespace))
	}

	// Get the kubeconfig.
	if name == "" {
		kcs, err := o.Client.ListKubeconfigs(o.Context, o.TCP)
		if err != nil {
			return errors.Wrap(err, "error listing kubeconfigs")
		}

		//raw, err := json.Marshal(kcs)
		raw, err := yaml.Marshal(kcs)
		if err != nil {
			return errors.Wrap(err, "error marshaling kubeconfigs")
		}

		if len(kcs) > 0 {
			switch o.OutputFormat {
			case "yaml":
				o.Printer.PrintRaw(string(raw))
			case "table":
				o.Printer.PrintTable(output.KubeconfigListTable(kcs...))
			default:
				o.Printer.PrintRaw(string(raw))
			}

			return nil
		}
		o.Printer.Warningf("no kubeconfigs found")

		return nil
	}

	kubeconfig, err := o.Client.GetKubeconfig(o.Context, o.TCP)
	if err != nil {
		return errors.Wrap(err, "error getting kubeconfig")
	}

	switch o.OutputFormat {
	case "yaml":
		o.Printer.PrintRaw(string(kubeconfig.Data))
	case "table":
		o.Printer.PrintTable(output.KubeconfigListTable(*kubeconfig))
	default:
		o.Printer.PrintRaw(string(kubeconfig.Data))
	}

	return nil
}
