package install

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"

	runclient "github.com/fluxcd/pkg/runtime/client"
	"sigs.k8s.io/cli-utils/pkg/object"
)

const (
	defaultVersion          = "2.2.0"
	commandName             = "install"
	commandShortDescription = "install the Kamaji Kubernetes operator"
	kustomization           = "kustomization.yaml"
	defaultDataStore        = "etcd"
	kamajiControllerManager = "kamaji"
	apps                    = "apps"
	deployment              = "Deployment"
	statefulset             = "StatefulSet"
)

var (
	components = []object.ObjMetadata{
		{
			Namespace: "",
			Name:      kamajiControllerManager,
			GroupKind: schema.GroupKind{Group: apps, Kind: deployment},
		},
		{
			Namespace: "",
			Name:      defaultDataStore,
			GroupKind: schema.GroupKind{Group: apps, Kind: statefulset},
		},
	}
	timeout           = 5 * time.Minute
	kubeclientOptions = new(runclient.Options)
)
