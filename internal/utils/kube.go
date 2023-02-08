package utils

import (
	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubeClient() (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error getting kube config")
	}

	kube, err := client.New(config, client.Options{
		Scheme: NewScheme(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting new kubernetes client")
	}

	return kube, nil
}

func NewKubeConfigWithOptions(rcg genericclioptions.RESTClientGetter) (*rest.Config, error) {
	config, err := rcg.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error getting kube config")
	}

	return config, nil
}

func NewScheme() *apiruntime.Scheme {
	scheme := apiruntime.NewScheme()
	addSchemes(scheme)

	return scheme
}

func addSchemes(scheme *apiruntime.Scheme) {
	kamajiv1alpha1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	apiextensionsv1.AddToScheme(scheme)
	admissionregistrationv1.AddToScheme(scheme)
	rbacv1.AddToScheme(scheme)
}
