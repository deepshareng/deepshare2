package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusDSUsage holds metrics for all dsusage methods.
type prometheusDSUsage struct {
	httpGetDuration       prometheus.Histogram
	httpDeleteDuration    prometheus.Histogram
	storageGetDuration    prometheus.Histogram
	storageDeleteDuration prometheus.Histogram
	storageIncDuration    prometheus.Histogram
}

var PrometheusForDSUsage = NewPrometheusForDSUsage()

func NewPrometheusForDSUsage() DSUsageInstrument {
	i := &prometheusDSUsage{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "dsusage",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		httpDeleteDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "dsusage",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP POST duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "dsusage",
				Name:      "storage_get_duration_milliseconds",
				Help:      "Bucketed histogram of storage get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageDeleteDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "dsusage",
				Name:      "storage_delete_duration_milliseconds",
				Help:      "Bucketed histogram of storage delete duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageIncDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "dsusage",
				Name:      "storage_inc_duration_milliseconds",
				Help:      "Bucketed histogram of storage inc duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.httpDeleteDuration)
	prometheus.MustRegister(i.storageGetDuration)
	prometheus.MustRegister(i.storageDeleteDuration)
	prometheus.MustRegister(i.storageIncDuration)

	return i
}

func (ps *prometheusDSUsage) HTTPGetDuration(start time.Time) {
	ps.httpGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
func (ps *prometheusDSUsage) HTTPDeleteDuration(start time.Time) {
	ps.httpDeleteDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusDSUsage) StorageGetDuration(start time.Time) {
	ps.storageGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusDSUsage) StorageDeleteDuration(start time.Time) {
	ps.storageDeleteDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusDSUsage) StorageIncDuration(start time.Time) {
	ps.storageIncDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
