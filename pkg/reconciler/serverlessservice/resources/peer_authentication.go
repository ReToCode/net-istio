package resources

import (
	"istio.io/client-go/pkg/apis/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"knative.dev/pkg/kmeta"

	securityv1beta1 "istio.io/api/security/v1beta1"
	typev1beta1 "istio.io/api/type/v1beta1"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
)

const (
	// TODO move those out to pkg or networking
	servingDrainPort = 8022
)

func MakeAllowDrainPeerAuthentication(sks *v1alpha1.ServerlessService, svcLister corev1listers.ServiceLister) (*v1beta1.PeerAuthentication, error) {
	svc, err := svcLister.Services(sks.Namespace).Get(sks.Status.PrivateServiceName)
	if err != nil {
		return nil, err
	}

	pa := &v1beta1.PeerAuthentication{
		ObjectMeta: metav1.ObjectMeta{
			Name:            sks.Status.PrivateServiceName,
			Namespace:       sks.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(sks)},
			Annotations:     sks.GetAnnotations(),
			Labels:          sks.GetLabels(),
		},
		Spec: securityv1beta1.PeerAuthentication{
			Selector: &typev1beta1.WorkloadSelector{
				MatchLabels: svc.Spec.Selector,
			},
			PortLevelMtls: map[uint32]*securityv1beta1.PeerAuthentication_MutualTLS{
				servingDrainPort: {
					Mode: securityv1beta1.PeerAuthentication_MutualTLS_PERMISSIVE,
				},
			},
		},
	}

	return pa, nil
}
