package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMatch holds metrics for all match methods.
type prometheusMatch struct {
	httpGetDuration       prometheus.Histogram
	httpPostDuration      prometheus.Histogram
	storageGetDuration    prometheus.Histogram
	storageSaveDuration   prometheus.Histogram
	storageDeleteDuration prometheus.Histogram
	storageHGetDuration   prometheus.Histogram
	storageHSetDuration   prometheus.Histogram
	storageHDelDuration   prometheus.Histogram
}

var PrometheusForMatch = NewPrometheusForMatch()

func NewPrometheusForMatch() Match {
	i := &prometheusMatch{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		httpPostDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP POST duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_get_duration_milliseconds",
				Help:      "Bucketed histogram of storage get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageSaveDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_save_duration_milliseconds",
				Help:      "Bucketed histogram of storage save duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageDeleteDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_delete_duration_milliseconds",
				Help:      "Bucketed histogram of storage delete duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageHSetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_hset_duration_milliseconds",
				Help:      "Bucketed histogram of storage delete duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			},
		),
		storageHGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_hget_duration_milliseconds",
				Help:      "Bucketed histogram of storage delete duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			},
		),
		storageHDelDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "match",
				Name:      "storage_hdel_duration_milliseconds",
				Help:      "Bucketed histogram of storage delete duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			},
		),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.httpPostDuration)
	prometheus.MustRegister(i.storageGetDuration)
	prometheus.MustRegister(i.storageSaveDuration)
	prometheus.MustRegister(i.storageDeleteDuration)
	prometheus.MustRegister(i.storageHSetDuration)
	prometheus.MustRegister(i.storageHGetDuration)
	prometheus.MustRegister(i.storageHDelDuration)

	return i
}

func (ps *prometheusMatch) HTTPGetDuration(start time.Time) {
	ps.httpGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
func (ps *prometheusMatch) HTTPPostDuration(start time.Time) {
	ps.httpPostDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageGetDuration(start time.Time) {
	ps.storageGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageSaveDuration(start time.Time) {
	ps.storageSaveDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageDeleteDuration(start time.Time) {
	ps.storageDeleteDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageHSetDuration(start time.Time) {
	ps.storageHSetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageHGetDuration(start time.Time) {
	ps.storageHGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusMatch) StorageHDelDuration(start time.Time) {
	ps.storageHDelDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
