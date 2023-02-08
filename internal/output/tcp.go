package output

import (
	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"

	"github.com/maxgio92/kamajictl/pkg/kamaji/tcp"
)

var (
	TCPListHeader        = []string{"NAMESPACE", "NAME", "VERSION", "STATUS", "ENDPOINT", "DATA STORE"}
	KubeconfigListHeader = []string{"NAMESPACE", "TCP", "NAME"}
)

func TCPListTable(tcps ...kamajiv1alpha1.TenantControlPlane) [][]string {
	var output [][]string

	output = append(output, TCPListHeader)

	for i := range tcps {
		var item []string

		item = append(item,
			tcps[i].Namespace,
			tcps[i].Name,
			tcps[i].Spec.Kubernetes.Version,
		)

		if tcps[i].Status.Kubernetes.Version.Status != nil {
			item = append(item,
				string(*tcps[i].Status.Kubernetes.Version.Status),
			)
		} else {
			item = append(item, "")
		}

		item = append(item,
			tcps[i].Status.ControlPlaneEndpoint,
			tcps[i].Status.Storage.DataStoreName,
		)

		output = append(output, item)
	}

	return output
}

func KubeconfigListTable(kubeconfigs ...tcp.KubeConfig) [][]string {
	var output [][]string

	output = append(output, KubeconfigListHeader)

	for i := range kubeconfigs {
		var item []string

		item = append(item,
			kubeconfigs[i].Metadata.Namespace,
			kubeconfigs[i].Metadata.OwnerTCP,
			kubeconfigs[i].Metadata.Name,
		)

		output = append(output, item)
	}

	return output
}
