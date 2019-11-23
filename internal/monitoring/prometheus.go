package monitoring

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// PrometheusService struct
type PrometheusService struct {
	Port string
}

// Serve Prometheus server
func (p *PrometheusService) Serve() {
	go func() {
		err := http.ListenAndServe(":"+p.Port, promhttp.Handler())
		if err != nil {
			log.Fatalf("Can't start monitoring server%s", err)
		}
	}()
}
