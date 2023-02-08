package tcp

import (
	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
)

// TCPOptions represent generic options of TenantControlPlanes
// for all the operations (create, update, get, list, delete).
type TCPOptions struct {
	*kamajiv1alpha1.TenantControlPlane
}

type TCPOption func(opts *TCPOptions)

func NewTCPOptions(opts ...TCPOption) *TCPOptions {
	o := &TCPOptions{&kamajiv1alpha1.TenantControlPlane{}}

	for _, f := range opts {
		f(o)
	}

	return o
}

func WithNamespace(namespace string) TCPOption {
	return func(opts *TCPOptions) {
		opts.Namespace = namespace
	}
}

func WithName(name string) TCPOption {
	return func(opts *TCPOptions) {
		opts.Name = name
	}
}

func WithServiceType(t string) TCPOption {
	return func(opts *TCPOptions) {
		opts.Spec.ControlPlane.Service.ServiceType = kamajiv1alpha1.ServiceType(t)
	}
}
func WithKubernetesVersion(v string) TCPOption {
	return func(opts *TCPOptions) {
		opts.Spec.Kubernetes.Version = v
	}
}

func (c *TCPOptions) Set(opts ...TCPOption) {
	for _, f := range opts {
		f(c)
	}
}

func (o *TCPOptions) Validate() error {
	if o.Name == "" {
		return ErrNameEmpty
	}

	if o.Namespace == "" {
		return ErrNamespaceEmpty
	}

	return nil
}

func (o *TCPOptions) validateList() error {
	if o.Namespace == "" {
		return ErrNamespaceEmpty
	}

	return nil
}

func (o *TCPOptions) AddCommonFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.Namespace, "namespace", "n", "",
		"The Namespace of the TenantControlPlane")
}

func (o *TCPOptions) AddCreateFlags(flags *pflag.FlagSet) {
	o.AddCommonFlags(flags)
	flags.StringVar(&o.Spec.Kubernetes.Version, "version", defaultKubernetesVersion,
		"The Kubernetes version of the TenantControlPlane")
	flags.StringVar((*string)(&o.Spec.ControlPlane.Service.ServiceType), "service-type",
		string(corev1.ServiceTypeClusterIP),
		"The Service type of the Tenant Control Plane (ClusterIP, NodePort or LB)")
}
