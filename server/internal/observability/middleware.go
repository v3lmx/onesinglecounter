package observability

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

func TestCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf(`requests_total{path=%q}`, r.URL.Path)
		metrics.GetOrCreateCounter(s).Inc()
		next.ServeHTTP(w, r)
	})
}
