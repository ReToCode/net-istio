/*
Copyright 2023 The Knative Authors

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

package resources

import (
	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/kmap"
	"knative.dev/pkg/kmeta"
)

const (
	// has to match config/700-istio-secret.yaml
	knativeServingCertsSecret = "knative-serving-certs"

	// has to match https://github.com/knative-sandbox/control-protocol/blob/main/pkg/certificates/constants.go#L21
	knativeFakeDnsName = "data-plane.knative.dev"
)

// MakeInternalEncryptionDestinationRule creates a DestinationRule that enables upstream TLS
// on for the specified host
func MakeInternalEncryptionDestinationRule(host string, ing *v1alpha1.Ingress, http2 bool) *v1alpha3.DestinationRule {
	dr := &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:            host,
			Namespace:       ing.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(ing)},
			Annotations:     ing.GetAnnotations(),
		},
		Spec: istiov1alpha3.DestinationRule{
			Host: host,
			TrafficPolicy: &istiov1alpha3.TrafficPolicy{
				Tls: &istiov1alpha3.ClientTLSSettings{
					Mode:            istiov1alpha3.ClientTLSSettings_SIMPLE,
					CredentialName:  knativeServingCertsSecret,
					SubjectAltNames: []string{knativeFakeDnsName},
				},
			},
		},
	}

	// Populate the Ingress labels.
	dr.Labels = kmap.Filter(ing.GetLabels(), func(k string) bool {
		return k != RouteLabelKey && k != RouteNamespaceLabelKey
	})
	dr.Labels[networking.IngressLabelKey] = ing.Name

	if http2 {
		dr.Spec.TrafficPolicy.ConnectionPool = &istiov1alpha3.ConnectionPoolSettings{
			Http: &istiov1alpha3.ConnectionPoolSettings_HTTPSettings{
				H2UpgradePolicy: istiov1alpha3.ConnectionPoolSettings_HTTPSettings_UPGRADE},
		}
	}

	return dr
}
