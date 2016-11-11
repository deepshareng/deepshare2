package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var PrometheusForInappData = newPrometheusInappData()

type prometheusInappData struct {
	httpPostDuration prometheus.Histogram
}

func newPrometheusInappData() InappData {
	i := &prometheusInappData{
		httpPostDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "inappdata",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP Post duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpPostDuration)

	return i
}

func (p *prometheusInappData) HTTPPostDuration(start time.Time) {
	p.httpPostDuration.Observe(float64(time.Since(start) / time.Millisecond))
}
