package config

import (
	"os"
	"path/filepath"

	"github.com/blakelead/nsinjector/pkg/generated/clientset/versioned"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")

// Clients struct regroups all necessary clients
type Clients struct {
	Kube      *kubernetes.Clientset
	Discovery *discovery.DiscoveryClient
	Dynamic   dynamic.Interface
	NRI       *versioned.Clientset
}

// NewClients creates all clients and return them in Clients struct
func NewClients() (*Clients, error) {
	config, err := clientcmd.BuildConfigFromFlags("", getConfig(kubeconfig))
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kubernetesClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	nriClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Discovery: discoveryClient,
		Dynamic:   dynamicClient,
		Kube:      kubernetesClient,
		NRI:       nriClient,
	}, nil
}

func getConfig(kubeconfig string) string {
	_, err := os.Stat(kubeconfig)
	if !os.IsNotExist(err) && err == nil {
		return kubeconfig
	}
	return ""
}
