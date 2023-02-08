package tcp

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

// KubeConfig represents a kubeconfig for a specific TenantControlPlane.
type KubeConfig struct {
	Metadata KubeConfigMetadata
	Data     []byte
}

// KubeConfigMetadata are KubeConfig metadata, including the reference to the
// TenantControlPlane that owns the kubeconfig that the metadata refers to.
type KubeConfigMetadata struct {
	Name      string
	Namespace string
	OwnerTCP  string
}

func NewKubeConfig(name, namespace, tcpowner string, content []byte) *KubeConfig {
	return &KubeConfig{
		Metadata: KubeConfigMetadata{
			Name:      name,
			Namespace: namespace,
			OwnerTCP:  tcpowner,
		},
		Data: content,
	}
}

func (m *KubeConfig) MarshalJSON() ([]byte, error) {
	name := `"name":"` + m.Metadata.Name + `"`
	namespace := `"namespace":"` + m.Metadata.Namespace + `"`
	tcpOwner := `"owner_tcp":"` + m.Metadata.OwnerTCP + `"`

	var config clientcmdapiv1.Config
	if err := yaml.Unmarshal(m.Data, &config); err != nil {
		return nil, err
	}

	content, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	ret := []byte(`{"metadata": {` + name + `,` + namespace + `,` + tcpOwner + `},"data":` + string(content) + `}`)

	return ret, nil
}
