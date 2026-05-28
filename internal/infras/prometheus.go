package infras

import (
	"context"

	"noahingh/hermes-agent-operator/internal/usecase"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type PrometheusMetric struct {
	reconcileTotal    *prometheus.CounterVec
	reconcileDuration prometheus.Histogram
	configMapOps      *prometheus.CounterVec
	statefulSetOps    *prometheus.CounterVec
	notFoundTotal     prometheus.Counter
	managed           prometheus.Gauge
}

func NewPrometheusMetric() *PrometheusMetric {
	m := &PrometheusMetric{
		reconcileTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hermesagent_reconcile_total",
			Help: "Total number of HermesAgent reconciliations.",
		}, []string{"result"}),
		reconcileDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "hermesagent_reconcile_duration_seconds",
			Help:    "Duration of HermesAgent reconciliations in seconds.",
			Buckets: prometheus.DefBuckets,
		}),
		configMapOps: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hermesagent_configmap_operations_total",
			Help: "Total number of ConfigMap create/update operations.",
		}, []string{"operation", "result"}),
		statefulSetOps: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hermesagent_statefulset_operations_total",
			Help: "Total number of StatefulSet create/update operations.",
		}, []string{"operation", "result"}),
		notFoundTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "hermesagent_not_found_total",
			Help: "Total number of reconciliations where the HermesAgent was not found.",
		}),
		managed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hermesagent_managed_total",
			Help: "Current number of HermesAgent resources managed by the operator.",
		}),
	}

	metrics.Registry.MustRegister(
		m.reconcileTotal,
		m.reconcileDuration,
		m.configMapOps,
		m.statefulSetOps,
		m.notFoundTotal,
		m.managed,
	)

	return m
}

func (m *PrometheusMetric) IncReconcile(_ context.Context, param usecase.IncReconcileParam) {
	m.reconcileTotal.WithLabelValues(param.Result.String()).Inc()
}

func (m *PrometheusMetric) ObserveReconcileDuration(_ context.Context, param usecase.ObserveReconcileDurationParam) {
	m.reconcileDuration.Observe(param.Seconds)
}

func (m *PrometheusMetric) IncConfigMapOperation(_ context.Context, param usecase.IncConfigMapOperationParam) {
	m.configMapOps.WithLabelValues(param.Operation.String(), param.Result.String()).Inc()
}

func (m *PrometheusMetric) IncStatefulSetOperation(_ context.Context, param usecase.IncStatefulSetOperationParam) {
	m.statefulSetOps.WithLabelValues(param.Operation.String(), param.Result.String()).Inc()
}

func (m *PrometheusMetric) IncNotFound(_ context.Context, _ usecase.IncNotFoundParam) {
	m.notFoundTotal.Inc()
}

func (m *PrometheusMetric) SetManaged(_ context.Context, param usecase.SetManagedParam) {
	m.managed.Set(float64(param.Count))
}
