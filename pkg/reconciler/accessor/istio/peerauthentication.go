/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package istio

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	istioclientset "knative.dev/net-istio/pkg/client/istio/clientset/versioned"
	istiolisters "knative.dev/net-istio/pkg/client/istio/listers/security/v1beta1"
	kaccessor "knative.dev/net-istio/pkg/reconciler/accessor"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
)

// PeerAuthenticationAccessor is an interface for accessing PeerAuthentication.
type PeerAuthenticationAccessor interface {
	GetIstioClient() istioclientset.Interface
	GetPeerAuthenticationLister() istiolisters.PeerAuthenticationLister
}

func isDifferent(current, desired *v1beta1.PeerAuthentication) bool {
	return !cmp.Equal(&current.Spec, &desired.Spec, protocmp.Transform()) ||
		!cmp.Equal(current.Labels, desired.Labels) ||
		!cmp.Equal(current.Annotations, desired.Annotations)
}

// ReconcilePeerAuthentication reconciles PeerAuthentication to the desired status.
func ReconcilePeerAuthentication(ctx context.Context, owner kmeta.Accessor, desired *v1beta1.PeerAuthentication,
	paAccessor PeerAuthenticationAccessor) (*v1beta1.PeerAuthentication, error) {

	recorder := controller.GetEventRecorder(ctx)
	if recorder == nil {
		return nil, fmt.Errorf("recorder for reconciling PeerAuthentication %s/%s is not created", desired.Namespace, desired.Name)
	}
	ns := desired.Namespace
	name := desired.Name
	dr, err := paAccessor.GetPeerAuthenticationLister().PeerAuthentications(ns).Get(name)
	if apierrs.IsNotFound(err) {
		dr, err = paAccessor.GetIstioClient().SecurityV1beta1().PeerAuthentications(ns).Create(ctx, desired, metav1.CreateOptions{})
		if err != nil {
			recorder.Eventf(owner, corev1.EventTypeWarning, "CreationFailed",
				"Failed to create PeerAuthentication %s/%s: %v", ns, name, err)
			return nil, fmt.Errorf("failed to create PeerAuthentication: %w", err)
		}
		recorder.Eventf(owner, corev1.EventTypeNormal, "Created", "Created PeerAuthentication %q", desired.Name)
	} else if err != nil {
		return nil, err
	} else if !metav1.IsControlledBy(dr, owner) {
		// Return an error with NotControlledBy information.
		return nil, kaccessor.NewAccessorError(
			fmt.Errorf("owner: %s with Type %T does not own PeerAuthentication: %q", owner.GetName(), owner, name),
			kaccessor.NotOwnResource)
	} else if isDifferent(dr, desired) {
		// Don't modify the informers copy
		existing := dr.DeepCopy()
		existing.Spec = *desired.Spec.DeepCopy()
		existing.Labels = desired.Labels
		existing.Annotations = desired.Annotations
		dr, err = paAccessor.GetIstioClient().SecurityV1beta1().PeerAuthentications(ns).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to update PeerAuthentication: %w", err)
		}
		recorder.Eventf(owner, corev1.EventTypeNormal, "Updated", "Updated PeerAuthentication %s/%s", ns, name)
	}
	return dr, nil
}
