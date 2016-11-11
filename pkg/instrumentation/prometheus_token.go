package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var PrometheusForToken = newPrometheusToken()

type prometheusToken struct {
	httpGetDuration prometheus.Histogram
}

func newPrometheusToken() Token {
	i := &prometheusToken{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "token",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP Get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)

	return i
}

func (p *prometheusToken) HTTPGetDuration(start time.Time) {
	p.httpGetDuration.Observe(float64(time.Since(start) / time.Millisecond))
}
