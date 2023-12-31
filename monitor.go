package main

import "github.com/prometheus/client_golang/prometheus"

// Monitor represents a Prometheus monitor
// It contains Prometheus registry and all available metrics
type Monitor struct {
	Registry            *prometheus.Registry
	ObjectDistanceGauge *prometheus.GaugeVec
}

// NewMonitor creates a new Monitor
func NewMonitor() *Monitor {
	reg := prometheus.NewRegistry()
	monitor := &Monitor{
		Registry: reg,

		ObjectDistanceGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "distance_between_objects",
			Help: "Current distance between objects in the space in meters",
		}, []string{"from", "to"}),
	}

	reg.MustRegister(monitor.ObjectDistanceGauge)

	return monitor
}
