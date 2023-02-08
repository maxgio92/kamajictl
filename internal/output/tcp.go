package output

import (
	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
)

var (
	TCPListHeader = []string{"NAME", "VERSION", "STATUS", "ENDPOINT", "DATA STORE"}
)

func TCPListTable(tcps ...kamajiv1alpha1.TenantControlPlane) [][]string {
	var output [][]string

	output = append(output, TCPListHeader)

	for i := range tcps {
		var item []string

		item = append(item,
			tcps[i].Name,
			tcps[i].Spec.Kubernetes.Version,
			string(*tcps[i].Status.Kubernetes.Version.Status),
			tcps[i].Status.ControlPlaneEndpoint,
			tcps[i].Status.Storage.DataStoreName,
		)

		output = append(output, item)
	}

	return output
}
