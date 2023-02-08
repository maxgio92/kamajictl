package options

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/maxgio92/kamajicli/internal/output"
)

// CommonOptions represents the output options of the command line.
type CommonOptions struct {
	Context   context.Context
	Verbosity int8
	Logger    *log.Logger
	Printer   *output.Printer
}

type CommonOption func(opts *CommonOptions)

func NewCommonOptions(opts ...CommonOption) *CommonOptions {
	o := &CommonOptions{}

	for _, f := range opts {
		f(o)
	}

	return o
}

func WithLogger(logger *log.Logger) CommonOption {
	return func(opts *CommonOptions) {
		opts.Logger = logger
	}
}

func WithContext(ctx context.Context) CommonOption {
	return func(opts *CommonOptions) {
		opts.Context = ctx
	}
}

func WithPrinter(printer *output.Printer) CommonOption {
	return func(opts *CommonOptions) {
		opts.Printer = printer
	}
}

func (o *CommonOptions) Set(opts ...CommonOption) {
	for _, f := range opts {
		f(o)
	}
}

func (o *CommonOptions) Validate() error {
	if o.Verbosity > 3 || o.Verbosity < 1 {
		return fmt.Errorf("verbosity level is not valid, valid numbers are between 1 and 3")
	}

	return nil
}

func (o *CommonOptions) AddFlags(flags *pflag.FlagSet) {
	flags.Int8VarP(&o.Verbosity, "verbosity", "v", 1, "Verbosity level")
}
