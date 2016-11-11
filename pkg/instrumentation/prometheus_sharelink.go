package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var PrometheusForSharelink = newPrometheusSharelink()

type prometheusSharelink struct {
	httpGetDuration prometheus.Histogram
}

func newPrometheusSharelink() Sharelink {
	i := &prometheusSharelink{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "sharelink",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP Get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)

	return i
}

func (p *prometheusSharelink) HTTPGetDuration(start time.Time) {
	p.httpGetDuration.Observe(float64(time.Since(start) / time.Millisecond))
}
