package options

import (
	"context"
	"embed"
	"fmt"
	"github.com/maxgio92/kamajictl/internal/output/log"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/maxgio92/kamajictl/internal/output"
)

// CommonOptions represents the output options of the command line.
type CommonOptions struct {
	Context           context.Context
	Verbosity         int8
	Printer           *output.Printer
	Logger            log.Logger
	Embedded          *embed.FS
	KubeconfigOptions *genericclioptions.ConfigFlags

	// Whether to get TenantControlPlanes across all namespaces.
	AllNamespaces bool
}

type CommonOption func(opts *CommonOptions)

func NewCommonOptions(opts ...CommonOption) *CommonOptions {
	o := &CommonOptions{
		KubeconfigOptions: genericclioptions.NewConfigFlags(true),
	}

	for _, f := range opts {
		f(o)
	}

	return o
}

func WithEmbeddedFS(embedded *embed.FS) CommonOption {
	return func(opts *CommonOptions) {
		opts.Embedded = embedded
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

func WithLogger(logger log.Logger) CommonOption {
	return func(opts *CommonOptions) {
		opts.Logger = logger
	}
}

func WithNamespace(namespace string) CommonOption {
	return func(opts *CommonOptions) {
		opts.KubeconfigOptions.Namespace = &namespace
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
	if o.KubeconfigOptions != nil {
		o.KubeconfigOptions.AddFlags(flags)
	}
}
