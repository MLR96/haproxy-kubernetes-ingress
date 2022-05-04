// Copyright 2019 HAProxy Technologies LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package customresource

import (
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/haproxytech/kubernetes-ingress/crs/generated/clientset/versioned"
	informers "github.com/haproxytech/kubernetes-ingress/crs/generated/informers/externalversions"
	"github.com/haproxytech/kubernetes-ingress/pkg/k8s"
	"github.com/haproxytech/kubernetes-ingress/pkg/store"
	"github.com/haproxytech/kubernetes-ingress/pkg/utils"
)

var logger = utils.GetLogger()

const (
	CoreGroupVersion = "core.haproxy.org/v1alpha1"
)

// A Custom Resource interface
// Any CR should be able to provide its kind, its kubernetes Informer
// and a method to process the update of a CR
type CR interface {
	GetKind() string
	GetInformer(chan k8s.SyncDataEvent, informers.SharedInformerFactory) cache.SharedIndexInformer
	ProcessEvent(*store.K8s, k8s.SyncDataEvent) bool
}

type CRManager struct {
	crs         map[string]CR
	client      *clientset.Clientset
	store       *store.K8s
	cacheResync time.Duration
	channel     chan k8s.SyncDataEvent
	stop        chan struct{}
}

func NewCRManager(s *store.K8s, restCfg *rest.Config, cacheResync time.Duration, eventChan chan k8s.SyncDataEvent, stop chan struct{}) CRManager {
	manager := CRManager{
		crs:         map[string]CR{},
		client:      clientset.NewForConfigOrDie(restCfg),
		store:       s,
		cacheResync: cacheResync,
		channel:     eventChan,
		stop:        stop,
	}
	manager.RegisterCoreCR(NewGlobalCR())
	manager.RegisterCoreCR(NewDefaultsCR())
	manager.RegisterCoreCR(NewBackendCR())
	return manager
}

func (m CRManager) RegisterCoreCR(cr CR) {
	resources, err := m.client.DiscoveryClient.ServerResourcesForGroupVersion(CoreGroupVersion)
	if err != nil {
		return
	}
	logger.Debugf("Custom API core.haproxy.org available")
	kindName := cr.GetKind()
	for _, resource := range resources.APIResources {
		if resource.Kind == kindName {
			m.crs[kindName] = cr
			logger.Infof("%s CR defined in API core.haproxy.org", kindName)
			break
		}
	}
}

func (m CRManager) EventCustomResource(job k8s.SyncDataEvent) bool {
	if cr, ok := m.crs[job.CRKind]; ok {
		return cr.ProcessEvent(m.store, job)
	}
	return false
}

// RunInformers runs Custom Resource Informers and return an array of corresponding cache.InformerSynced
func (m CRManager) RunInformers(namespace string) []cache.InformerSynced {
	informerSynced := make([]cache.InformerSynced, 0, len(m.crs))
	informerFactory := informers.NewSharedInformerFactoryWithOptions(m.client, m.cacheResync, informers.WithNamespace(namespace))
	for _, controller := range m.crs {
		informer := controller.GetInformer(m.channel, informerFactory)
		go informer.Run(m.stop)
		informerSynced = append(informerSynced, informer.HasSynced)
	}
	return informerSynced
}
