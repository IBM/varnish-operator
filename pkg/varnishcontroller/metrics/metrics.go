package metrics

import "github.com/prometheus/client_golang/prometheus"

const vclCompilationErrorMetricName = "varnish_vcl_compilation_error"

type VarnishControllerMetrics struct {
	VCLCompilationError prometheus.Gauge
}

func NewVarnishControllerMetrics() *VarnishControllerMetrics {
	return &VarnishControllerMetrics{
		VCLCompilationError: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: vclCompilationErrorMetricName,
				Help: "Indicates if the VCL compilation failed. 0 - successfully compiled, 1 - failed.",
			},
		),
	}
}
