package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCheckAuthCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_check_auth_count",
		Help: "API check auth",
	})

	apiCheckAuthBlockedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_check_auth_blocked_count",
		Help: "API check auth blocked",
	})
)

func init() {
	prometheus.MustRegister(apiCheckAuthCounter)
	prometheus.MustRegister(apiCheckAuthBlockedCounter)
}
