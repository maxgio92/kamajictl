package kamaji

import (
	"context"
	"github.com/maxgio92/kamajictl/internal/output/log"

	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/maxgio92/kamajictl/pkg/kamaji/tcp"
)

// KamajiClient represents the options of the command line for managing Tenant Control Plane resources.
type KamajiClient struct {
	kube   client.Client
	logger log.Logger
	*tcp.TCPClient
}

type KamajiClientOption func(opts *KamajiClient)

func NewKamajiClient(opts ...KamajiClientOption) (*KamajiClient, error) {
	o := &KamajiClient{}

	for _, f := range opts {
		f(o)
	}

	if err := o.Validate(); err != nil {
		return nil, err
	}

	to, err := tcp.NewTCPClient(
		tcp.WithLogger(o.logger),
		tcp.WithKubeClient(o.kube),
	)
	if err != nil {
		return nil, err
	}

	o.TCPClient = to

	return o, nil
}

func WithLogger(logger log.Logger) KamajiClientOption {
	return func(opts *KamajiClient) {
		opts.logger = logger
	}
}

func WithKubeClient(kube client.Client) KamajiClientOption {
	return func(opts *KamajiClient) {
		opts.kube = kube
	}
}

func (c *KamajiClient) Set(opts ...KamajiClientOption) {
	for _, f := range opts {
		f(c)
	}
}

func (c *KamajiClient) Validate() error {
	if c.kube == nil {
		return ErrKubeClientNil
	}

	if c.logger == nil {
		return ErrLoggerNil
	}

	return nil
}

// PurgeCRDs deletes the TenantControlPlane CRD.
// ATTENTION: please ensure all the TenantControlPlane resources are deleted by calling PurgeTCPs before calling this.
func (c *KamajiClient) PurgeCRDs(ctx context.Context) error {
	var list apiextensionsv1.CustomResourceDefinitionList

	selector := client.MatchingLabels{resourcesSelectorKey: resourcesSelectorValue}

	if err := c.kube.List(ctx, &list, selector); err != nil {
		return errors.Wrap(err, "error listing Kamaji CRDs")
	}

	for _, r := range list.Items {
		if err := c.kube.Delete(ctx, &r, &client.DeleteOptions{}); err != nil {
			c.logger.Warningf("error deleting %s CRD: %s", r.Name, err)
		}
		c.logger.Infof("CRD %s deleted", r.Name)
	}

	return nil
}

// PurgeWebhooks deletes the Kamaji validating and mutating webhook configurations.
func (c *KamajiClient) PurgeWebhooks(ctx context.Context) error {
	var vlist admissionregistrationv1.ValidatingWebhookConfigurationList

	selector := client.MatchingLabels{resourcesSelectorKey: resourcesSelectorValue}

	if err := c.kube.List(ctx, &vlist, selector); err != nil {
		return errors.Wrap(err, "error listing kamaji validating webhook configurations")
	}

	for _, r := range vlist.Items {
		if err := c.kube.Delete(ctx, &r, &client.DeleteOptions{}); err != nil {
			c.logger.Warningf("error deleting validatingwebhookconfiguration %s: %s", r.Name, err)
		}
		c.logger.Infof("validatingwebhookconfiguration %s deleted", r.Name)
	}

	var mlist admissionregistrationv1.MutatingWebhookConfigurationList

	if err := c.kube.List(ctx, &mlist, selector); err != nil {
		return errors.Wrap(err, "error listing kamaji validating webhook configurations")
	}

	for _, r := range mlist.Items {
		if err := c.kube.Delete(ctx, &r, &client.DeleteOptions{}); err != nil {
			c.logger.Warningf("error deleting mutatingwebhookconfiguration %s: %s", r.Name, err)
		}
		c.logger.Infof("mutatingwebhookconfiguration %s deleted", r.Name)
	}

	return nil
}

func (c *KamajiClient) PurgeRBAC(ctx context.Context) error {
	var err error
	if err = c.purgeClusterRoles(ctx); err != nil {
		return err
	}

	if err = c.purgeClusterRoleBindings(ctx); err != nil {
		return err
	}

	return nil
}

func (c *KamajiClient) purgeClusterRoles(ctx context.Context) error {
	var list rbacv1.ClusterRoleList

	selector := client.MatchingLabels{resourcesSelectorKey: resourcesSelectorValue}

	if err := c.kube.List(ctx, &list, selector); err != nil {
		return errors.Wrap(err, "error listing kamaji clusterroles")
	}

	for _, r := range list.Items {
		if err := c.kube.Delete(ctx, &r, &client.DeleteOptions{}); err != nil {
			c.logger.Warningf("error deleting clusterrole %s: %s", r.Name, err)
		}
		c.logger.Infof("clusterrole %s deleted", r.Name)
	}

	return nil
}

func (c *KamajiClient) purgeClusterRoleBindings(ctx context.Context) error {
	var list rbacv1.ClusterRoleBindingList

	selector := client.MatchingLabels{resourcesSelectorKey: resourcesSelectorValue}

	if err := c.kube.List(ctx, &list, selector); err != nil {
		return errors.Wrap(err, "error listing kamaji clusterrolebindings")
	}

	for _, r := range list.Items {
		if err := c.kube.Delete(ctx, &r, &client.DeleteOptions{}); err != nil {
			c.logger.Warningf("error deleting clusterrolebinding %s: %s", r.Name, err)
		}
		c.logger.Infof("clusterrolebinding %s deleted", r.Name)
	}

	return nil
}
