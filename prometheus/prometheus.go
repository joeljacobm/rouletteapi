package prometheus

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Hits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "rouletteapi",
			Name:      "page_hit_counter",
			Help:      `Total numbe of hits for each URL path,partitioned by response code`,
		},
		[]string{"path"},
	)

	ErrorCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "rouletteapi",
			Name:      "errors_total",
			Help:      `Total number of errors`,
		},
	)

	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rouletteapi_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(HttpDuration.WithLabelValues(path))
		Hits.With(prometheus.Labels{"path": path}).Inc()
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}
