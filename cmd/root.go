package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/maxgio92/kamajicli/cmd/create"
	"github.com/maxgio92/kamajicli/cmd/delete"
	"github.com/maxgio92/kamajicli/cmd/get"
	"github.com/maxgio92/kamajicli/internal/options"
	"github.com/maxgio92/kamajicli/internal/output"
)

var (
	out = bufio.NewWriter(os.Stdout)
)

// NewRootCommand returns the root command.
func NewRootCommand(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	// Add operation sub-commands.
	cmd.AddCommand(create.NewCreateCmd(opts))
	cmd.AddCommand(delete.NewDeleteCmd(opts))
	cmd.AddCommand(get.NewGetCmd(opts))

	// Add common options persistent flags.
	opts.AddFlags(cmd.PersistentFlags())

	return cmd
}

// Execute creates the root command, adds all child commands, and execute the root command.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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

	// Instantiate the logger.
	logger := log.New()

	// Instantiate the printer.
	printer := output.NewPrinter(
		output.WithWriter(out),
	)

	// Instantiate the common command-line options.
	opts := options.NewCommonOptions(
		options.WithContext(ctx),
		options.WithLogger(logger),
		options.WithPrinter(printer),
	)

	// Create the root command, propagate the context to subcommands, and execute it.
	opts.Printer.CheckErr(NewRootCommand(opts).Execute())
}
