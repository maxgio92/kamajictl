package cmd

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"github.com/maxgio92/kamajictl/cmd/install"
	"github.com/maxgio92/kamajictl/cmd/uninstall"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajictl/cmd/create"
	"github.com/maxgio92/kamajictl/cmd/delete"
	"github.com/maxgio92/kamajictl/cmd/get"
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/output"
)

var (
	out = bufio.NewWriter(os.Stdout)
)

type Options struct {
	*options.CommonOptions
}

// NewRootCommand returns the root command.
func NewRootCommand(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:               commandName,
		Short:             commandShortDescription,
		DisableAutoGenTag: true,
	}

	// Add operation sub-commands.
	cmd.AddCommand(create.NewCreateCmd(opts))
	cmd.AddCommand(delete.NewDeleteCmd(opts))
	cmd.AddCommand(get.NewGetCmd(opts))
	cmd.AddCommand(install.NewInstallCommand(opts))
	cmd.AddCommand(uninstall.NewUninstallCommand(opts))

	// Add common options persistent flags.
	opts.AddFlags(cmd.PersistentFlags())

	return cmd
}

// Execute creates the root command, adds all child commands, and execute the root command.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(embedded *embed.FS) {
	// Ensure the buffer is flushed to output when returning.
	defer out.Flush()

	// Mark context done when signals arrive.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Gracefully shut down and reset signals behaviour.
	go func() {
		<-ctx.Done()
		fmt.Println("Terminating...")
		stop()
	}()

	// Instantiate the common command-line options.
	opts := options.NewCommonOptions(
		options.WithContext(ctx),
		options.WithLogger(output.NewPrinter(output.WithWriter(os.Stderr))),
		options.WithPrinter(output.NewPrinter(output.WithWriter(out))),
		options.WithLogger(output.NewPrinter(output.WithWriter(os.Stderr))),
		options.WithEmbeddedFS(embedded),
		options.WithNamespace("kamaji-system"),
	)

	// Create the root command, propagate the context to subcommands, and execute it.
	opts.Printer.CheckErr(NewRootCommand(opts).Execute())
}
