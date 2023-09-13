package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var reconcilerDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "hugopage_reconciler_duration",
		Help: "How long the reconcile loop ran for in microseconds",
	},
	[]string{
		"reconciler",
		"name",
		"namespace",
	},
)

var active = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "hugopage_active",
		Help: "Number of installed hugopage objects for this instance",
	},
	[]string{
		"type",
	},
)

func init() {
	metrics.Registry.MustRegister(reconcilerDuration)
	metrics.Registry.MustRegister(active)
}
