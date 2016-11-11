package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMatch holds metrics for all match methods.
type prometheusBindDeviceToCookie struct {
	httpGetDuration prometheus.Histogram
}

var PrometheusForBindDeviceToCookier = NewPrometheusForBindDeviceToCookier()

func NewPrometheusForBindDeviceToCookier() BindDeviceToCookierInstrument {
	i := &prometheusBindDeviceToCookie{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "frontend_bindDeviceToCookier",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	return i
}

func (ps *prometheusBindDeviceToCookie) HTTPGetDuration(start time.Time) {
	ps.httpGetDuration.Observe(float64(time.Since(start)) / float64(time.Millisecond))
}
