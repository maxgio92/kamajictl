package get

import (
	"github.com/maxgio92/kamajictl/cmd/get/kubeconfig"
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajictl/cmd/get/tcp"
)

// NewGetCmd returns the Delete command.
func NewGetCmd(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	cmd.PersistentFlags().BoolVarP(&opts.AllNamespaces, "all-namespaces", "A", false, "Get object(s) across all namespaces")

	cmd.AddCommand(tcp.NewTCPCommand(opts))
	cmd.AddCommand(kubeconfig.NewKubeconfigCommand(opts))

	return cmd
}
