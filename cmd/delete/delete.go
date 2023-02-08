package delete

import (
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajictl/cmd/delete/tcp"
)

// NewDeleteCmd returns the Delete command.
func NewDeleteCmd(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	cmd.AddCommand(tcp.NewTCPCommand(opts))

	return cmd
}
