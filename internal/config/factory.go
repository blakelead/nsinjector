package config

import (
	"time"

	"github.com/blakelead/nsinjector/pkg/generated/informers/externalversions"
	"k8s.io/client-go/informers"
)

// InformerFactories regroups all necessary informer factories
type InformerFactories struct {
	Kube informers.SharedInformerFactory
	NRI  externalversions.SharedInformerFactory
}

// NewFactories creates and returns all factories
func NewFactories(clients *Clients, defaultResync time.Duration) *InformerFactories {
	return &InformerFactories{
		Kube: informers.NewSharedInformerFactory(clients.Kube, defaultResync),
		NRI:  externalversions.NewSharedInformerFactory(clients.NRI, defaultResync),
	}
}

// Start starts all informer factories
func (f *InformerFactories) Start(stop <-chan struct{}) {
	f.Kube.Start(stop)
	f.NRI.Start(stop)
}
