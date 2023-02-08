package create

import (
	"github.com/maxgio92/kamajicli/internal/options"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajicli/cmd/create/tcp"
)

// NewCreateCmd returns the Create command.
func NewCreateCmd(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	cmd.AddCommand(tcp.NewTCPCommand(opts))

	return cmd
}
