package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zebradil/airgradient-exporter/pkg/collector"
	"github.com/zebradil/airgradient-exporter/pkg/logger"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	fmt.Printf("airgradient-exporter version %s (commit: %s, date: %s)\n", version, commit, date)
}

func printHelp() {
	fmt.Println(`AirGradient Exporter - A Prometheus exporter for AirGradient monitors

Usage:
  airgradient-exporter [flags]

Flags:
  --help     Show this help message
  --version  Show version information

Environment Variables:
  AIRGRADIENT_MONITORS    Comma-separated list of AirGradient monitor IPs, hostnames, or host:port (required)
  PORT                    Port to listen on (default: 9112)
  LOG_FORMAT              Log format: "text" or "json" (default: "text", colored when output is a terminal)

Examples:
  export AIRGRADIENT_MONITORS="192.168.1.50,192.168.1.51"
  airgradient-exporter

  export AIRGRADIENT_MONITORS="192.168.1.50:8080,192.168.1.51:8080"
  airgradient-exporter

  export AIRGRADIENT_MONITORS="192.168.1.50"
  export PORT="8080"
  airgradient-exporter

Visit http://localhost:9112/metrics to see the metrics.`)
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	// Get log format from environment variable
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}

	// Create logger with configurable format
	log := logger.NewLogger(logFormat, os.Stdout)

	monitorsEnv := os.Getenv("AIRGRADIENT_MONITORS")
	if monitorsEnv == "" {
		log.Error("AIRGRADIENT_MONITORS environment variable is required")
		os.Exit(1)
	}
	hosts := strings.Split(monitorsEnv, ",")
	for i := range hosts {
		hosts[i] = strings.TrimSpace(hosts[i])
	}

	log.Info("Starting airgradient-exporter", "monitors", hosts)

	// Create a new registry.
	reg := prometheus.NewRegistry()

	// Add standard process and go metrics
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	// Create and register our collector
	c := collector.NewAirGradientCollector(log, hosts)
	reg.MustRegister(c)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`<html>
			<head><title>AirGradient Exporter</title></head>
			<body>
			<h1>AirGradient Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`)); err != nil {
			log.Error("Failed to write response", "error", err)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9112"
	}

	addr := ":" + port
	log.Info("Listening on", "address", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
