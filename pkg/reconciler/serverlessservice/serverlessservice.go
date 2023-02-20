/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package serverlessservice

import (
	"context"
	"errors"
	"fmt"

	istioclientset "knative.dev/net-istio/pkg/client/istio/clientset/versioned"
	istiolisters "knative.dev/net-istio/pkg/client/istio/listers/networking/v1alpha3"
	"knative.dev/net-istio/pkg/reconciler/ingress/config"
	sksreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/serverlessservice"

	istioaccessor "knative.dev/net-istio/pkg/reconciler/accessor/istio"
	"knative.dev/net-istio/pkg/reconciler/serverlessservice/resources"
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	pkgreconciler "knative.dev/pkg/reconciler"
)

// reconciler implements controller.Reconciler for SKS resources.
type reconciler struct {
	istioclient istioclientset.Interface

	virtualServiceLister  istiolisters.VirtualServiceLister
	destinationRuleLister istiolisters.DestinationRuleLister
}

// Check that our Reconciler implements various interfaces.
var (
	_ sksreconciler.Interface               = (*reconciler)(nil)
	_ istioaccessor.VirtualServiceAccessor  = (*reconciler)(nil)
	_ istioaccessor.DestinationRuleAccessor = (*reconciler)(nil)
)

// Reconcile compares the actual state with the desired, and attempts to converge the two.
func (r *reconciler) ReconcileKind(ctx context.Context, sks *netv1alpha1.ServerlessService) pkgreconciler.Event {
	networkCfg := config.FromContext(ctx).Network
	if !networkCfg.EnableMeshPodAddressability && !networkCfg.InternalEncryption {
		// Just ignore if we're disabled.
		return nil
	}

	if networkCfg.EnableMeshPodAddressability && networkCfg.InternalEncryption {
		return errors.New("failed to reconcile VirtualService. EnableMeshPodAddressability and InternalEncryption cannot both be enabled")
	}

	if sks.Status.PrivateServiceName == "" {
		// No private service yet, nothing to do here.
		return nil
	}

	if networkCfg.EnableMeshPodAddressability {
		vs := resources.MakeVirtualService(sks)
		if _, err := istioaccessor.ReconcileVirtualService(ctx, sks, vs, r); err != nil {
			return fmt.Errorf("failed to reconcile VirtualService: %w", err)
		}

		dr := resources.MakeMeshAddressableDestinationRule(sks)
		if _, err := istioaccessor.ReconcileDestinationRule(ctx, sks, dr, r); err != nil {
			return fmt.Errorf("failed to reconcile DestinationRule: %w", err)
		}
	}

	if networkCfg.InternalEncryption {
		dr := resources.MakeInternalEncryptionDestinationRule(sks)
		if _, err := istioaccessor.ReconcileDestinationRule(ctx, sks, dr, r); err != nil {
			return fmt.Errorf("failed to reconcile DestinationRule: %w", err)
		}
	}

	return nil
}

func (r *reconciler) GetIstioClient() istioclientset.Interface {
	return r.istioclient
}

func (r *reconciler) GetVirtualServiceLister() istiolisters.VirtualServiceLister {
	return r.virtualServiceLister
}

func (r *reconciler) GetDestinationRuleLister() istiolisters.DestinationRuleLister {
	return r.destinationRuleLister
}
