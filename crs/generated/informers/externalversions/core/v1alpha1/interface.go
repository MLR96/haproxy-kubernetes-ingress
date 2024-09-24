//
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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/haproxytech/kubernetes-ingress/crs/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Backends returns a BackendInformer.
	Backends() BackendInformer
	// Defaults returns a DefaultsInformer.
	Defaults() DefaultsInformer
	// Globals returns a GlobalInformer.
	Globals() GlobalInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Backends returns a BackendInformer.
func (v *version) Backends() BackendInformer {
	return &backendInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Defaults returns a DefaultsInformer.
func (v *version) Defaults() DefaultsInformer {
	return &defaultsInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Globals returns a GlobalInformer.
func (v *version) Globals() GlobalInformer {
	return &globalInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}