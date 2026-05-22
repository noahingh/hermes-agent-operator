package usecase

import (
	"context"
	agentsv1alpha1 "noahingh/hermes-agent-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type HermesAgentUseCase struct {
	kube Kubernetes
}

type ReconcileParam struct {
	NamespacedName types.NamespacedName
}

func NewHermesAgentUseCase(kube Kubernetes) *HermesAgentUseCase {
	return &HermesAgentUseCase{
		kube,
	}
}

func (u *HermesAgentUseCase) Reconcile(ctx context.Context, param ReconcileParam) error {
	ha, err := u.kube.GetHermesAgent(ctx, GetHermesAgentParam{
		NamespacedName: param.NamespacedName,
	})
	if err != nil {
		return err
	}
	if ha == nil {
		return nil
	}

	sts, err := u.kube.GetStatefulSet(ctx, GetStatefulSetParam{
		NamespacedName: param.NamespacedName,
	})
	if err != nil {
		return err
	}

	desired := u.buildStatefulSet(ha)

	if sts != nil {
		desired.ResourceVersion = sts.ResourceVersion
		return u.kube.UpdateStatefulSet(ctx, UpdateStatefulSetParam{StatefulSet: desired})
	}

	return u.kube.CreateStatefulSetOwnedByHermesAgent(ctx, CreateStatefulSetOfHermesAgentParam{
		HermesAgent: ha,
		StatefulSet: desired,
	})
}

func (u *HermesAgentUseCase) buildStatefulSet(ha *agentsv1alpha1.HermesAgent) *appsv1.StatefulSet {
	replicas := int32(1)
	dshmSize := resource.MustParse("1Gi")

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ha.Name,
			Namespace: ha.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": ha.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": ha.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:             "hermes-agent",
							Image:            "nousresearch/hermes-agent:latest",
							ImagePullPolicy:  corev1.PullIfNotPresent,
							Args:             []string{"gateway", "run"},
							WorkingDir:       "/opt/hermes",
							Env: []corev1.EnvVar{
								{Name: "HERMES_HOME", Value: "/opt/data"},
								{Name: "HOME", Value: "/opt/data/home"},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("2"),
									corev1.ResourceMemory: resource.MustParse("4Gi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "dshm", MountPath: "/dev/shm"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "dshm",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium:    corev1.StorageMediumMemory,
									SizeLimit: &dshmSize,
								},
							},
						},
					},
				},
			},
		},
	}
}

