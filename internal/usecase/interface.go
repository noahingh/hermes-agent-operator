package usecase

import (
	"context"
	agentsv1alpha1 "noahingh/hermes-agent-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Kubernetes interface {
	GetHermesAgent(ctx context.Context, param GetHermesAgentParam) (*agentsv1alpha1.HermesAgent, error)

	GetConfigMap(ctx context.Context, param GetConfigMapParam) (*corev1.ConfigMap, error)
	CreateConfigMapOwnedByHermesAgent(ctx context.Context, param CreateConfigMapOfHermesAgentParam) error
	UpdateConfigMapOwnedByHermesAgent(ctx context.Context, param UpdateConfigMapParam) error

	GetStatefulSet(ctx context.Context, param GetStatefulSetParam) (*appsv1.StatefulSet, error)
	CreateStatefulSetOwnedByHermesAgent(ctx context.Context, param CreateStatefulSetOfHermesAgentParam) error
	UpdateStatefulSetOwnedByHermesAgent(ctx context.Context, param UpdateStatefulSetParam) error

	GetServiceAccount(ctx context.Context, param GetServiceAccountParam) (*corev1.ServiceAccount, error)
	CreateServiceAccountOwnedByHermesAgent(ctx context.Context, param CreateServiceAccountOfHermesAgentParam) error
	UpdateServiceAccountOwnedByHermesAgent(ctx context.Context, param UpdateServiceAccountParam) error
	DeleteServiceAccount(ctx context.Context, param DeleteServiceAccountParam) error

	GetRole(ctx context.Context, param GetRoleParam) (*rbacv1.Role, error)
	CreateRoleOwnedByHermesAgent(ctx context.Context, param CreateRoleOfHermesAgentParam) error
	UpdateRoleOwnedByHermesAgent(ctx context.Context, param UpdateRoleParam) error
	DeleteRole(ctx context.Context, param DeleteRoleParam) error

	GetRoleBinding(ctx context.Context, param GetRoleBindingParam) (*rbacv1.RoleBinding, error)
	CreateRoleBindingOwnedByHermesAgent(ctx context.Context, param CreateRoleBindingOfHermesAgentParam) error
	UpdateRoleBindingOwnedByHermesAgent(ctx context.Context, param UpdateRoleBindingParam) error
	DeleteRoleBinding(ctx context.Context, param DeleteRoleBindingParam) error
}

type GetHermesAgentParam struct {
	NamespacedName types.NamespacedName
}

type GetConfigMapParam struct {
	NamespacedName types.NamespacedName
}

type CreateConfigMapOfHermesAgentParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	ConfigMap   *corev1.ConfigMap
}

type UpdateConfigMapParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	ConfigMap   *corev1.ConfigMap
}

type GetStatefulSetParam struct {
	NamespacedName types.NamespacedName
}

type CreateStatefulSetOfHermesAgentParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	StatefulSet *appsv1.StatefulSet
}

type UpdateStatefulSetParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	StatefulSet *appsv1.StatefulSet
}

type GetServiceAccountParam struct {
	NamespacedName types.NamespacedName
}

type CreateServiceAccountOfHermesAgentParam struct {
	HermesAgent    *agentsv1alpha1.HermesAgent
	ServiceAccount *corev1.ServiceAccount
}

type UpdateServiceAccountParam struct {
	HermesAgent    *agentsv1alpha1.HermesAgent
	ServiceAccount *corev1.ServiceAccount
}

type DeleteServiceAccountParam struct {
	NamespacedName types.NamespacedName
}

type GetRoleParam struct {
	NamespacedName types.NamespacedName
}

type CreateRoleOfHermesAgentParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	Role        *rbacv1.Role
}

type UpdateRoleParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	Role        *rbacv1.Role
}

type DeleteRoleParam struct {
	NamespacedName types.NamespacedName
}

type GetRoleBindingParam struct {
	NamespacedName types.NamespacedName
}

type CreateRoleBindingOfHermesAgentParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	RoleBinding *rbacv1.RoleBinding
}

type UpdateRoleBindingParam struct {
	HermesAgent *agentsv1alpha1.HermesAgent
	RoleBinding *rbacv1.RoleBinding
}

type DeleteRoleBindingParam struct {
	NamespacedName types.NamespacedName
}
