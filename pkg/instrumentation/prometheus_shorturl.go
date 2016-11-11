package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusShorturl holds metrics for all shorturl methods.
type prometheusShorturl struct {
	httpGetDuration     prometheus.Histogram
	httpPostDuration    prometheus.Histogram
	storageGetDuration  prometheus.Histogram
	storageSaveDuration prometheus.Histogram
}

var PrometheusForShorturl = NewPrometheusForShorturl()

func NewPrometheusForShorturl() Shorturl {
	i := &prometheusShorturl{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "shorturl",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		httpPostDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "shorturl",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP POST duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "shorturl",
				Name:      "storage_get_duration_milliseconds",
				Help:      "Bucketed histogram of storage get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageSaveDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "shorturl",
				Name:      "storage_save_duration_milliseconds",
				Help:      "Bucketed histogram of storage save duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.httpPostDuration)
	prometheus.MustRegister(i.storageGetDuration)
	prometheus.MustRegister(i.storageSaveDuration)

	return i
}

func (ps *prometheusShorturl) HTTPGetDuration(start time.Time) {
	ps.httpGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
func (ps *prometheusShorturl) HTTPPostDuration(start time.Time) {
	ps.httpPostDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusShorturl) StorageGetDuration(start time.Time) {
	ps.storageGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusShorturl) StorageSaveDuration(start time.Time) {
	ps.storageSaveDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
