package monitoring

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type PrometheusService struct {
	Port string
}

func (p *PrometheusService) Serve() {
	go func() {
		err := http.ListenAndServe(":"+p.Port, promhttp.Handler())
		if err != nil {
			log.Fatalf("Can't start monitoring server%s", err)
		}
	}()
}
