package tcp

import (
	"context"
	"fmt"
	"github.com/maxgio92/kamajictl/internal/output/log"

	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TCPClient represents the options of the command line for managing Tenant Control Plane resources.
type TCPClient struct {
	kube   client.Client
	logger log.Logger
}

type TCPClientOption func(opts *TCPClient)

func NewTCPClient(opts ...TCPClientOption) (*TCPClient, error) {
	o := &TCPClient{}

	for _, f := range opts {
		f(o)
	}

	if err := o.Validate(); err != nil {
		return nil, err
	}

	return o, nil
}

func WithLogger(logger log.Logger) TCPClientOption {
	return func(opts *TCPClient) {
		opts.logger = logger
	}
}

func WithKubeClient(kube client.Client) TCPClientOption {
	return func(opts *TCPClient) {
		opts.kube = kube
	}
}

func (c *TCPClient) Set(opts ...TCPClientOption) {
	for _, f := range opts {
		f(c)
	}
}

func (c *TCPClient) Validate() error {
	if c.kube == nil {
		return ErrKubeClientNil
	}

	if c.logger == nil {
		return ErrLoggerNil
	}

	return nil
}

// CreateTCP creates TCPOptions resource using the provided context,
// and returns an error.
func (c *TCPClient) CreateTCP(ctx context.Context, opts *TCPOptions) error {
	if err := c.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}

	if err := opts.Validate(); err != nil {
		return errors.Wrap(err, "error validating create options")
	}

	tcp := &kamajiv1alpha1.TenantControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
		},
		Spec: kamajiv1alpha1.TenantControlPlaneSpec{
			Kubernetes: kamajiv1alpha1.KubernetesSpec{
				Version: opts.Spec.Kubernetes.Version,
			},
			ControlPlane: kamajiv1alpha1.ControlPlane{
				Deployment: kamajiv1alpha1.DeploymentSpec{},
				Service: kamajiv1alpha1.ServiceSpec{
					kamajiv1alpha1.AdditionalMetadata{},
					opts.Spec.ControlPlane.Service.ServiceType,
				},
			},
		},
	}

	if err := c.kube.Create(ctx, tcp, &client.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

// DeleteTCP deletes a TCPOptions resource using the provided context, and returns an error.
func (c *TCPClient) DeleteTCP(ctx context.Context, opts *TCPOptions) error {
	if err := c.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}

	if err := opts.Validate(); err != nil {
		return errors.Wrap(err, "error validating options")
	}

	tcp := &kamajiv1alpha1.TenantControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
		},
	}

	if err := c.kube.Delete(ctx, tcp); err != nil {
		return err
	}

	return nil
}

// GetTCP get a TCPOptions resource using the provided context, and returns an error.
func (c *TCPClient) GetTCP(ctx context.Context, opts *TCPOptions) (*kamajiv1alpha1.TenantControlPlane, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating options")
	}

	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating options")
	}

	tcp := &kamajiv1alpha1.TenantControlPlane{}

	key := types.NamespacedName{
		Namespace: opts.Namespace,
		Name:      opts.Name,
	}

	if err := c.kube.Get(ctx, key, tcp); err != nil {
		return nil, err
	}

	return tcp, nil
}

// ListTCPs get a TCPOptions resource using the provided context, and returns an error.
func (c *TCPClient) ListTCPs(ctx context.Context, opts *TCPOptions) (*kamajiv1alpha1.TenantControlPlaneList, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating client options")
	}

	tcpList := &kamajiv1alpha1.TenantControlPlaneList{}

	listOptions := &client.ListOptions{}
	if opts.Namespace != "" {
		listOptions = &client.ListOptions{Namespace: opts.Namespace}
	}

	if err := c.kube.List(ctx, tcpList, listOptions); err != nil {
		return nil, err
	}

	return tcpList, nil
}

// GetKubeconfig retrieves a kubeconfig for a specified TenantControlPlane using the provided context, and returns the
// kubeconfig as string and an error.
func (c *TCPClient) GetKubeconfig(ctx context.Context, opts *TCPOptions) (*KubeConfig, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating client options")
	}

	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating options")
	}

	t, err := c.GetTCP(ctx, opts)
	if err != nil {
		return nil, err
	}

	var secretName string
	if t.Status.KubeConfig.Admin.SecretName == "" {
		return nil, fmt.Errorf("kubeconfig secret for tenantcontrolplane %s/%s not ready", opts.Name, opts.Namespace)
	}

	secretName = t.Status.KubeConfig.Admin.SecretName

	key := types.NamespacedName{
		Namespace: opts.Namespace,
		Name:      secretName,
	}

	secret := &corev1.Secret{}

	if err := c.kube.Get(ctx, key, secret); err != nil {
		return nil, err
	}

	return NewKubeConfig(key.Name, key.Namespace, t.Name, secret.Data["admin.conf"]), nil
}

// ListKubeconfigs retrieves a kubeconfig for a specified TenantControlPlane using the provided context, and returns the
// kubeconfig as string and an error.
func (c *TCPClient) ListKubeconfigs(ctx context.Context, opts *TCPOptions) ([]KubeConfig, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "error validating client options")
	}

	ts, err := c.ListTCPs(ctx, opts)
	if err != nil {
		return nil, err
	}

	if ts.Items == nil {
		return nil, errors.New("error getting tenantcontrolplanes")
	}

	kcs := []KubeConfig{}

	for _, t := range ts.Items {
		t := t

		var secretName string
		if t.Status.KubeConfig.Admin.SecretName == "" {
			return nil, fmt.Errorf("kubeconfig secret for tenantcontrolplane %s/%s not ready", opts.Name, opts.Namespace)
		}

		secretName = t.Status.KubeConfig.Admin.SecretName

		key := types.NamespacedName{
			Namespace: t.Namespace,
			Name:      secretName,
		}

		secret := &corev1.Secret{}

		if err := c.kube.Get(ctx, key, secret); err != nil {
			return nil, err
		}

		kcs = append(kcs, *NewKubeConfig(key.Name, key.Namespace, t.Name, secret.Data["admin.conf"]))
	}

	return kcs, nil
}

// PurgeTCPs deletes all the TenantControlPlane at cluster-level.
// ATTENTION: it deletes all the resources across all the Namespaces.
func (c *TCPClient) PurgeTCPs(ctx context.Context) error {
	list, err := c.ListTCPs(ctx, NewTCPOptions())
	if err != nil {
		return err
	}

	for _, r := range list.Items {
		do := NewTCPOptions(
			WithName(r.Name),
			WithNamespace(r.Namespace),
		)
		if err := c.DeleteTCP(ctx, do); err != nil {
			c.logger.Warningf("error deleting %s/%s TCP: %s", r.Namespace, r.Name, err)
		}
		c.logger.Infof("TCP %s/%s deleted", r.Name, r.Namespace)
	}

	return nil
}
