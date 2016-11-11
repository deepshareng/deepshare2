package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMatch holds metrics for all match methods.
type prometheusAppCookieDevice struct {
	httpGetDuration    prometheus.Histogram
	httpPutDuration    prometheus.Histogram
	storageGetDuration prometheus.Histogram
	storagePutDuration prometheus.Histogram
}

var PrometheusForAppCookieDevice = NewPrometheusForAppCookieDevice()

func NewPrometheusForAppCookieDevice() AppCookieDeviceInstrument {
	i := &prometheusAppCookieDevice{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appCookieDevice",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		httpPutDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appCookieDevice",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP POST duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appCookieDevice",
				Name:      "storage_get_duration_milliseconds",
				Help:      "Bucketed histogram of storage get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storagePutDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appCookieDevice",
				Name:      "storage_save_duration_milliseconds",
				Help:      "Bucketed histogram of storage save duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.httpPutDuration)
	prometheus.MustRegister(i.storageGetDuration)
	prometheus.MustRegister(i.storagePutDuration)
	return i
}

func (ps *prometheusAppCookieDevice) HTTPGetDuration(start time.Time) {
	ps.httpGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusAppCookieDevice) HTTPPutDuration(start time.Time) {
	ps.httpPutDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusAppCookieDevice) StorageGetDuration(start time.Time) {
	ps.storageGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}

func (ps *prometheusAppCookieDevice) StoragePutDuration(start time.Time) {
	ps.storagePutDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
