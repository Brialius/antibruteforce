package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCheckAuthCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_check_auth_count",
		Help: "API check auth",
	}, []string{"ip", "login"})

	apiCheckAuthBlockedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_check_auth_blocked_count",
		Help: "API check auth blocked",
	}, []string{"ip", "login"})
)

func init() {
	prometheus.MustRegister(apiCheckAuthCounter)
	prometheus.MustRegister(apiCheckAuthBlockedCounter)
}
