package tcp_test

import (
	"context"
	"fmt"
	"github.com/maxgio92/kamajictl/internal/output/log"

	kamajiv1alpha1 "github.com/clastix/kamaji/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/maxgio92/kamajictl/internal/output"
	"github.com/maxgio92/kamajictl/pkg/kamaji/tcp"
)

var (
	existentTCPNames   = []string{"foo", "bar"}
	existentNamespaces = []string{"foo", "bar"}
)

var _ = Describe("A new TCP client", func() {
	var (
		logger     log.Logger
		kubeClient client.Client
		tcpClient  *tcp.TCPClient
		err        error
	)

	BeforeEach(func() {
		kubeClient = fakeKubeClient()
		Expect(kubeClient).ToNot(BeNil())

		logger = output.NewPrinter()
		Expect(logger).ToNot(BeNil())
	})

	When("logger and kube client are set", func() {
		BeforeEach(func() {
			tcpClient, err = tcp.NewTCPClient(tcp.WithLogger(logger), tcp.WithKubeClient(kubeClient))
		})
		It("should be built successfully", func() {
			Expect(tcpClient).ToNot(BeNil())
		})
		It("should not error", func() {
			Expect(err).To(BeNil())
		})
	})

	When("missing logger", func() {
		BeforeEach(func() {
			tcpClient, err = tcp.NewTCPClient(tcp.WithKubeClient(kubeClient))
		})
		It("should not be built successfully", func() {
			Expect(tcpClient).To(BeNil())
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	When("missing kube client", func() {
		BeforeEach(func() {
			tcpClient, err = tcp.NewTCPClient(tcp.WithLogger(logger))
		})
		It("should not be built successfully", func() {
			Expect(tcpClient).To(BeNil())
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	When("is not nil", func() {
		var tcpClient *tcp.TCPClient
		BeforeEach(func() {
			tcpClient, err = tcp.NewTCPClient(tcp.WithLogger(logger), tcp.WithKubeClient(kubeClient))
			Expect(tcpClient).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
		It("should be valid", func() {
			Expect(tcpClient.Validate()).To(BeNil())
		})
	})
})

var _ = Describe("Getting a TCP", func() {
	var (
		logger          log.Logger
		kubeClient      client.Client
		tcpClient       *tcp.TCPClient
		name, namespace string
		err             error
	)

	BeforeEach(func() {
		name = sampleTCPList().Items[0].Name
		namespace = sampleTCPList().Items[0].Namespace

		kubeClient = fakeKubeClient()
		Expect(kubeClient).ToNot(BeNil())

		logger = output.NewPrinter()
		Expect(logger).ToNot(BeNil())

		tcpClient, err = tcp.NewTCPClient(
			tcp.WithLogger(logger),
			tcp.WithKubeClient(kubeClient))
		Expect(tcpClient).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	When("passing only the name", func() {
		var opts *tcp.TCPOptions
		var tcps *kamajiv1alpha1.TenantControlPlane
		var err error
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithName(name))
			tcps, err = tcpClient.GetTCP(context.Background(), opts)
		})
		It("should fail", func() {
			Expect(err).ToNot(BeNil())
			Expect(tcps).To(BeNil())
		})
	})

	When("passing only the namespace", func() {
		var opts *tcp.TCPOptions
		var tcps *kamajiv1alpha1.TenantControlPlane
		var err error
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithNamespace(namespace))
			tcps, err = tcpClient.GetTCP(context.Background(), opts)
		})
		It("should fail", func() {
			Expect(err).ToNot(BeNil())
			Expect(tcps).To(BeNil())
		})
	})

	When("passing name and namespace", func() {
		var (
			opts *tcp.TCPOptions
			tcps *kamajiv1alpha1.TenantControlPlane
			err  error
		)

		When("there is a match", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(name),
					tcp.WithNamespace(namespace))
				tcps, err = tcpClient.GetTCP(context.Background(), opts)
			})
			It("should not error", func() {
				Expect(err).To(BeNil())
			})
			It("should give a result", func() {
				Expect(tcps).ToNot(BeNil())
			})
			It("the result should match the query", func() {
				Expect(tcps.ObjectMeta.Name).To(Equal(name))
				Expect(tcps.ObjectMeta.Namespace).To(Equal(namespace))
			})
		})

		When("there is not a match", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(existentTCPNames[0]),
					tcp.WithNamespace(existentNamespaces[1]))
				tcps, err = tcpClient.GetTCP(context.Background(), opts)
			})
			It("should error", func() {
				Expect(err).ToNot(BeNil())
			})
			It("should not give a result", func() {
				Expect(tcps).To(BeNil())
			})
		})
	})
})

var _ = Describe("Listing TCPs", func() {
	var (
		logger     log.Logger
		kubeClient client.Client
		tcpClient  *tcp.TCPClient
		opts       *tcp.TCPOptions
		namespace  string
		err        error
	)

	BeforeEach(func() {
		namespace = sampleTCPList().Items[0].Namespace

		kubeClient = fakeKubeClient()
		Expect(kubeClient).ToNot(BeNil())

		logger = output.NewPrinter()
		Expect(logger).ToNot(BeNil())

		tcpClient, err = tcp.NewTCPClient(
			tcp.WithLogger(logger),
			tcp.WithKubeClient(kubeClient))
		Expect(tcpClient).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	Context("with the namespace", func() {
		var tcps *kamajiv1alpha1.TenantControlPlaneList
		var err error

		When("there is a match", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(tcp.WithNamespace(namespace))
				tcps, err = tcpClient.ListTCPs(context.Background(), opts)
			})
			It("should not error", func() {
				Expect(err).To(BeNil())
			})
			It("should give a result", func() {
				Expect(tcps.Items).ToNot(BeNil())
				Expect(len(tcps.Items)).To(BeNumerically(">=", 1))
			})
		})

		When("there is not a match", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(tcp.WithNamespace("notexists"))
				tcps, err = tcpClient.ListTCPs(context.Background(), opts)
			})
			It("should not error", func() {
				Expect(err).To(BeNil())
			})
			It("should not give a result", func() {
				Expect(tcps.Items).ToNot(BeNil())
				Expect(len(tcps.Items)).To(BeNumerically("==", 0))
			})
		})
	})
	Context("without the namespace", func() {
		var tcps *kamajiv1alpha1.TenantControlPlaneList
		var err error
		BeforeEach(func() {
			opts = tcp.NewTCPOptions()
			tcps, err = tcpClient.ListTCPs(context.Background(), opts)
		})
		It("should not error", func() {
			Expect(err).To(BeNil())
		})
		It("should give a result", func() {
			Expect(tcps).ToNot(BeNil())
			Expect(len(tcps.Items)).To(BeNumerically("==", 2))
		})
	})
})

var _ = Describe("Creating a TCP", func() {
	var (
		logger     log.Logger
		kubeClient client.Client
		tcpClient  *tcp.TCPClient
		opts       *tcp.TCPOptions
		err        error
	)

	BeforeEach(func() {
		kubeClient = fakeKubeClient()
		logger = output.NewPrinter()
		tcpClient, err = tcp.NewTCPClient(
			tcp.WithLogger(logger),
			tcp.WithKubeClient(kubeClient))
		Expect(tcpClient).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	Context("with only the TCP name", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithName(existentTCPNames[0]))
			err = tcpClient.CreateTCP(context.Background(), opts)
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with only the TCP namespace", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithNamespace(existentNamespaces[0]))
			err = tcpClient.CreateTCP(context.Background(), opts)
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with TCP name and namespace", func() {
		When("a TCP does not exist with same name and namespace", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(existentTCPNames[0]),
					tcp.WithNamespace(existentNamespaces[1]),
				)
				err = tcpClient.CreateTCP(context.Background(), opts)
			})
			It("should not error", func() {
				Expect(err).To(BeNil())
			})
		})

		When("a TCP already exists with same name and namespace", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(existentTCPNames[0]),
					tcp.WithNamespace(existentNamespaces[0]),
				)
				err = tcpClient.CreateTCP(context.Background(), opts)
			})
			It("should error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("Deleting a TCP", func() {
	var (
		logger     log.Logger
		kubeClient client.Client
		tcpClient  *tcp.TCPClient
		opts       *tcp.TCPOptions
		err        error
	)

	BeforeEach(func() {
		kubeClient = fakeKubeClient()
		logger = output.NewPrinter()
		tcpClient, err = tcp.NewTCPClient(
			tcp.WithLogger(logger),
			tcp.WithKubeClient(kubeClient))
		Expect(tcpClient).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	Context("with only the TCP name", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithName(existentTCPNames[0]))
			err = tcpClient.DeleteTCP(context.Background(), opts)
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with only the TCP namespace", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithNamespace(existentNamespaces[0]))
			err = tcpClient.DeleteTCP(context.Background(), opts)
		})
		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with TCP name and namespace", func() {
		When("a TCP does not exist with specified name and namespace", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(existentTCPNames[0]),
					tcp.WithNamespace(existentNamespaces[1]),
				)
				err = tcpClient.DeleteTCP(context.Background(), opts)
			})
			It("should error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("a TCP exists with specified name and namespace", func() {
			BeforeEach(func() {
				opts = tcp.NewTCPOptions(
					tcp.WithName(existentTCPNames[0]),
					tcp.WithNamespace(existentNamespaces[0]),
				)
				err = tcpClient.DeleteTCP(context.Background(), opts)
			})
			It("should not error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})

var _ = Describe("Getting a kubeconfig", func() {
	var (
		logger          log.Logger
		kubeClient      client.Client
		tcpClient       *tcp.TCPClient
		name, namespace string
		kubeconfig      *tcp.KubeConfig
		err             error
		opts            *tcp.TCPOptions
	)

	BeforeEach(func() {
		name = sampleTCPList().Items[0].Name
		namespace = sampleTCPList().Items[0].Namespace

		kubeClient = fakeKubeClient()
		Expect(kubeClient).ToNot(BeNil())

		logger = output.NewPrinter()
		Expect(logger).ToNot(BeNil())

		tcpClient, err = tcp.NewTCPClient(
			tcp.WithLogger(logger),
			tcp.WithKubeClient(kubeClient))
		Expect(tcpClient).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	When("passing only the name", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithName(name))
			kubeconfig, err = tcpClient.GetKubeconfig(context.Background(), opts)
		})
		It("should fail", func() {
			Expect(err).ToNot(BeNil())
			Expect(kubeconfig).To(BeNil())
		})
	})

	When("passing only the namespace", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithNamespace(namespace))
			kubeconfig, err = tcpClient.GetKubeconfig(context.Background(), opts)
		})
		It("should fail", func() {
			Expect(err).ToNot(BeNil())
			Expect(kubeconfig).To(BeNil())
		})
	})

	When("passiong name and namespace", func() {
		BeforeEach(func() {
			opts = tcp.NewTCPOptions(tcp.WithName(name), tcp.WithNamespace(namespace))
			kubeconfig, err = tcpClient.GetKubeconfig(context.Background(), opts)
		})
		It("should work", func() {
			Expect(err).To(BeNil())
			Expect(kubeconfig.Data).ToNot(BeEmpty())
		})
	})

})

func fakeKubeClient() client.Client {
	kubeClientBuilder := fake.NewClientBuilder()

	scheme := apiruntime.NewScheme()
	kamajiv1alpha1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)

	kubeClientBuilder.WithScheme(scheme)
	kubeClientBuilder.WithLists(sampleTCPList())
	kubeClientBuilder.WithObjects(sampleKubeconfigSecret())

	return kubeClientBuilder.Build()
}

func sampleKubeconfigSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-admin-kubeconfig", sampleTCPList().Items[0].Name),
			Namespace: sampleTCPList().Items[0].Namespace,
		},
		Data: map[string][]byte{
			"admin.conf": []byte("fake kubeconfig"),
		},
	}
}

func sampleTCPList() *kamajiv1alpha1.TenantControlPlaneList {
	return &kamajiv1alpha1.TenantControlPlaneList{
		Items: []kamajiv1alpha1.TenantControlPlane{
			{
				metav1.TypeMeta{},
				metav1.ObjectMeta{
					Name:      existentTCPNames[0],
					Namespace: existentNamespaces[0],
				},
				kamajiv1alpha1.TenantControlPlaneSpec{
					Kubernetes: kamajiv1alpha1.KubernetesSpec{
						Version:              "1.26.0",
						Kubelet:              kamajiv1alpha1.KubeletSpec{},
						AdmissionControllers: nil,
					},
					ControlPlane: kamajiv1alpha1.ControlPlane{
						Deployment: kamajiv1alpha1.DeploymentSpec{
							Replicas: 2,
						},
						Service: kamajiv1alpha1.ServiceSpec{},
						Ingress: nil,
					},
				},
				kamajiv1alpha1.TenantControlPlaneStatus{
					KubeConfig: kamajiv1alpha1.KubeconfigsStatus{
						Admin: kamajiv1alpha1.KubeconfigStatus{
							SecretName: fmt.Sprintf("%s-admin-kubeconfig", existentTCPNames[0]),
						},
					},
				},
			},
			{
				metav1.TypeMeta{},
				metav1.ObjectMeta{
					Name:      existentTCPNames[1],
					Namespace: existentNamespaces[1],
				},
				kamajiv1alpha1.TenantControlPlaneSpec{
					Kubernetes: kamajiv1alpha1.KubernetesSpec{
						Version:              "1.26.0",
						Kubelet:              kamajiv1alpha1.KubeletSpec{},
						AdmissionControllers: nil,
					},
					ControlPlane: kamajiv1alpha1.ControlPlane{
						Deployment: kamajiv1alpha1.DeploymentSpec{
							Replicas: 2,
						},
						Service: kamajiv1alpha1.ServiceSpec{},
						Ingress: nil,
					},
				},
				kamajiv1alpha1.TenantControlPlaneStatus{
					KubeConfig: kamajiv1alpha1.KubeconfigsStatus{
						Admin: kamajiv1alpha1.KubeconfigStatus{
							SecretName: fmt.Sprintf("%s-admin-kubeconfig", existentTCPNames[1]),
						},
					},
				},
			},
		},
	}
}
