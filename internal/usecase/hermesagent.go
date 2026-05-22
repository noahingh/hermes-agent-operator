package usecase

import (
	"context"
	"fmt"
	agentsv1alpha1 "noahingh/hermes-agent-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	sigsyaml "sigs.k8s.io/yaml"
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

	if err := u.reconcileConfigMap(ctx, ha); err != nil {
		return err
	}

	return u.reconcileStatefulSet(ctx, ha)
}

func (u *HermesAgentUseCase) reconcileConfigMap(ctx context.Context, ha *agentsv1alpha1.HermesAgent) error {
	cmName := configMapName(ha)
	cm, err := u.kube.GetConfigMap(ctx, GetConfigMapParam{
		NamespacedName: types.NamespacedName{Name: cmName, Namespace: ha.Namespace},
	})
	if err != nil {
		return err
	}

	desired, err := u.buildConfigMap(ha)
	if err != nil {
		return err
	}

	if cm != nil {
		desired.ResourceVersion = cm.ResourceVersion
		return u.kube.UpdateConfigMap(ctx, UpdateConfigMapParam{ConfigMap: desired})
	}

	return u.kube.CreateConfigMapOwnedByHermesAgent(ctx, CreateConfigMapOfHermesAgentParam{
		HermesAgent: ha,
		ConfigMap:   desired,
	})
}

func (u *HermesAgentUseCase) buildConfigMap(ha *agentsv1alpha1.HermesAgent) (*corev1.ConfigMap, error) {
	data := map[string]string{}
	if ha.Spec.Hermes != nil && ha.Spec.Hermes.Config != nil {
		yamlBytes, err := sigsyaml.JSONToYAML(ha.Spec.Hermes.Config.Raw)
		if err != nil {
			return nil, fmt.Errorf("converting config to YAML: %w", err)
		}
		data["config.yaml"] = string(yamlBytes)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(ha),
			Namespace: ha.Namespace,
		},
		Data: data,
	}, nil
}


func (u *HermesAgentUseCase) reconcileStatefulSet(ctx context.Context, ha *agentsv1alpha1.HermesAgent) error {
	sts, err := u.kube.GetStatefulSet(ctx, GetStatefulSetParam{
		NamespacedName: types.NamespacedName{Name: ha.Name, Namespace: ha.Namespace},
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
	hasConfig := ha.Spec.Hermes != nil && ha.Spec.Hermes.Config != nil

	volumes := []corev1.Volume{
		{
			Name: "dshm",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium:    corev1.StorageMediumMemory,
					SizeLimit: &dshmSize,
				},
			},
		},
	}

	mainVolumeMounts := []corev1.VolumeMount{
		{Name: "dshm", MountPath: "/dev/shm"},
	}

	var initContainers []corev1.Container

	if hasConfig {
		volumes = append(volumes,
			corev1.Volume{
				Name:         "data",
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
			corev1.Volume{
				Name: "bootstrap",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: configMapName(ha)},
					},
				},
			},
		)

		mainVolumeMounts = append(mainVolumeMounts,
			corev1.VolumeMount{Name: "data", MountPath: "/opt/data"},
		)

		initContainers = []corev1.Container{
			{
				Name:            "bootstrap-config",
				Image:           "nousresearch/hermes-agent:latest",
				ImagePullPolicy: corev1.PullIfNotPresent,
				Command:         []string{"/bin/sh", "-ec"},
				Args: []string{`set -eu
mkdir -p "/opt/data/home"
if [ -f "/bootstrap/config.yaml" ]; then
  cp "/bootstrap/config.yaml" "/opt/data/config.yaml"
fi
`},
				Env: []corev1.EnvVar{
					{Name: "HERMES_HOME", Value: "/opt/data"},
				},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "data", MountPath: "/opt/data"},
					{Name: "bootstrap", MountPath: "/bootstrap", ReadOnly: true},
				},
			},
		}
	}

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
					InitContainers: initContainers,
					Containers: []corev1.Container{
						{
							Name:            "hermes-agent",
							Image:           "nousresearch/hermes-agent:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args:            []string{"gateway", "run"},
							WorkingDir:      "/opt/hermes",
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
							VolumeMounts: mainVolumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}

func configMapName(ha *agentsv1alpha1.HermesAgent) string {
	return ha.Name + "-config"
}
