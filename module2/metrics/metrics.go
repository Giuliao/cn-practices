package metrics

import (
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const MetricsNamespace = "httpserver"

var (
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spent")
)

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		glog.Error(err)
	}
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"step"},
	)
}

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}

func NewTimer() *ExecutionTimer {
	return NewTimeExecutionTimer(functionLatency)
}

func NewTimeExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("totals").Observe(time.Now().Sub(t.start).Seconds())
}
