package infras

import (
	"context"

	"noahingh/hermes-agent-operator/internal/usecase"

	"github.com/prometheus/client_golang/prometheus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const telemetryLoggerName = "hermesagent"

type counterMetric struct {
	vec        *prometheus.CounterVec
	labelNames []string
}

type histogramMetric struct {
	vec        *prometheus.HistogramVec
	labelNames []string
}

type gaugeMetric struct {
	vec        *prometheus.GaugeVec
	labelNames []string
}

type PrometheusTelemetry struct {
	counters   map[string]counterMetric
	histograms map[string]histogramMetric
	gauges     map[string]gaugeMetric
}

func NewPrometheusTelemetry() *PrometheusTelemetry {
	t := &PrometheusTelemetry{
		counters:   map[string]counterMetric{},
		histograms: map[string]histogramMetric{},
		gauges:     map[string]gaugeMetric{},
	}

	t.registerCounter(usecase.MetricReconcileTotal,
		"Total number of HermesAgent reconciliations.",
		[]string{"result"})
	t.registerCounter(usecase.MetricConfigMapOperationsTotal,
		"Total number of ConfigMap create/update operations.",
		[]string{"operation", "result"})
	t.registerCounter(usecase.MetricStatefulSetOpsTotal,
		"Total number of StatefulSet create/update operations.",
		[]string{"operation", "result"})
	t.registerCounter(usecase.MetricNotFoundTotal,
		"Total number of reconciliations where the HermesAgent was not found.",
		nil)

	t.registerHistogram(usecase.MetricReconcileDurationSeconds,
		"Duration of HermesAgent reconciliations in seconds.",
		nil, prometheus.DefBuckets)

	t.registerGauge(usecase.MetricManagedTotal,
		"Current number of HermesAgent resources managed by the operator.",
		nil)

	return t
}

func (t *PrometheusTelemetry) registerCounter(name, help string, labels []string) {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: help}, labels)
	metrics.Registry.MustRegister(vec)
	t.counters[name] = counterMetric{vec: vec, labelNames: labels}
}

func (t *PrometheusTelemetry) registerHistogram(name, help string, labels []string, buckets []float64) {
	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: name, Help: help, Buckets: buckets}, labels)
	metrics.Registry.MustRegister(vec)
	t.histograms[name] = histogramMetric{vec: vec, labelNames: labels}
}

func (t *PrometheusTelemetry) registerGauge(name, help string, labels []string) {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labels)
	metrics.Registry.MustRegister(vec)
	t.gauges[name] = gaugeMetric{vec: vec, labelNames: labels}
}

func (t *PrometheusTelemetry) Info(ctx context.Context, msg string, keysAndValues ...any) {
	log.FromContext(ctx).WithName(telemetryLoggerName).Info(msg, keysAndValues...)
}

func (t *PrometheusTelemetry) Error(ctx context.Context, err error, msg string, keysAndValues ...any) {
	log.FromContext(ctx).WithName(telemetryLoggerName).Error(err, msg, keysAndValues...)
}

func (t *PrometheusTelemetry) IncCounter(name string, labels map[string]string) {
	m, ok := t.counters[name]
	if !ok {
		ctrl.Log.WithName(telemetryLoggerName).Info("unknown counter metric", "name", name)
		return
	}
	values, err := resolveLabelValues(m.labelNames, labels)
	if err != nil {
		ctrl.Log.WithName(telemetryLoggerName).Error(err, "invalid counter labels", "name", name)
		return
	}
	m.vec.WithLabelValues(values...).Inc()
}

func (t *PrometheusTelemetry) ObserveHistogram(name string, value float64, labels map[string]string) {
	m, ok := t.histograms[name]
	if !ok {
		ctrl.Log.WithName(telemetryLoggerName).Info("unknown histogram metric", "name", name)
		return
	}
	values, err := resolveLabelValues(m.labelNames, labels)
	if err != nil {
		ctrl.Log.WithName(telemetryLoggerName).Error(err, "invalid histogram labels", "name", name)
		return
	}
	m.vec.WithLabelValues(values...).Observe(value)
}

func (t *PrometheusTelemetry) SetGauge(name string, value float64, labels map[string]string) {
	m, ok := t.gauges[name]
	if !ok {
		ctrl.Log.WithName(telemetryLoggerName).Info("unknown gauge metric", "name", name)
		return
	}
	values, err := resolveLabelValues(m.labelNames, labels)
	if err != nil {
		ctrl.Log.WithName(telemetryLoggerName).Error(err, "invalid gauge labels", "name", name)
		return
	}
	m.vec.WithLabelValues(values...).Set(value)
}

func resolveLabelValues(names []string, labels map[string]string) ([]string, error) {
	if len(names) == 0 {
		return nil, nil
	}
	values := make([]string, len(names))
	for i, n := range names {
		v, ok := labels[n]
		if !ok {
			return nil, &missingLabelError{label: n}
		}
		values[i] = v
	}
	return values, nil
}

type missingLabelError struct {
	label string
}

func (e *missingLabelError) Error() string {
	return "missing label: " + e.label
}
