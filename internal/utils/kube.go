package utils

import (
	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubeClient() (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error getting kube config")
	}

	scheme := apiruntime.NewScheme()
	addSchemes(scheme)

	kube, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting new kubernetes client")
	}

	return kube, nil
}

func addSchemes(scheme *apiruntime.Scheme) {
	kamajiv1alpha1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
}
