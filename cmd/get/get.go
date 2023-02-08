package get

import (
	"github.com/maxgio92/kamajicli/cmd/get/kubeconfig"
	"github.com/maxgio92/kamajicli/internal/options"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajicli/cmd/get/tcp"
)

// NewGetCmd returns the Delete command.
func NewGetCmd(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	cmd.AddCommand(tcp.NewTCPCommand(opts))
	cmd.AddCommand(kubeconfig.NewKubeconfigCommand(opts))

	return cmd
}
