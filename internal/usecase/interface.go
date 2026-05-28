package usecase

import (
	"context"
	agentsv1alpha1 "noahingh/hermes-agent-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	MetricReconcileTotal           = "hermesagent_reconcile_total"
	MetricReconcileDurationSeconds = "hermesagent_reconcile_duration_seconds"
	MetricConfigMapOperationsTotal = "hermesagent_configmap_operations_total"
	MetricStatefulSetOpsTotal      = "hermesagent_statefulset_operations_total"
	MetricNotFoundTotal            = "hermesagent_not_found_total"
	MetricManagedTotal             = "hermesagent_managed_total"

	ResultSuccess  = "success"
	ResultError    = "error"
	ResultNotFound = "not_found"
	OpCreate       = "create"
	OpUpdate       = "update"
)

// Telemetry collects logs and metrics emitted by the usecase.
// Metric methods are generic and keyed by name; the implementation pre-registers
// each name with its collector type and label set. Adding a new metric only
// requires registering it in the implementation — no interface change.
type Telemetry interface {
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, err error, msg string, keysAndValues ...any)

	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
}

type Kubernetes interface {
	GetHermesAgent(ctx context.Context, param GetHermesAgentParam) (*agentsv1alpha1.HermesAgent, error)
	ListHermesAgents(ctx context.Context) ([]agentsv1alpha1.HermesAgent, error)

	GetConfigMap(ctx context.Context, param GetConfigMapParam) (*corev1.ConfigMap, error)
	CreateConfigMapOwnedByHermesAgent(ctx context.Context, param CreateConfigMapOfHermesAgentParam) error
	UpdateConfigMapOwnedByHermesAgent(ctx context.Context, param UpdateConfigMapParam) error

	GetStatefulSet(ctx context.Context, param GetStatefulSetParam) (*appsv1.StatefulSet, error)
	CreateStatefulSetOwnedByHermesAgent(ctx context.Context, param CreateStatefulSetOfHermesAgentParam) error
	UpdateStatefulSetOwnedByHermesAgent(ctx context.Context, param UpdateStatefulSetParam) error
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
