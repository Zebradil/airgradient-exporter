package collector

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/zebradil/airgradient-exporter/pkg/airgradient"
)

func TestNewAirGradientCollector(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	hosts := []string{"192.168.1.50"}

	collector := NewAirGradientCollector(logger, hosts)
	if collector == nil {
		t.Fatal("NewAirGradientCollector() returned nil")
	}
	if collector.client == nil {
		t.Fatal("client is nil")
	}
	if len(collector.hosts) != 1 || collector.hosts[0] != "192.168.1.50" {
		t.Errorf("hosts = %v, want [192.168.1.50]", collector.hosts)
	}
}

func TestAirGradientCollector_Describe(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	collector := NewAirGradientCollector(logger, []string{"test"})

	ch := make(chan *prometheus.Desc, 100)
	collector.Describe(ch)
	close(ch)

	descs := make([]*prometheus.Desc, 0)
	for desc := range ch {
		descs = append(descs, desc)
	}

	// Should have all metric descriptors: 1 (up) + 25 (measures) + 3 (config) = 29
	expectedCount := 29
	if len(descs) != expectedCount {
		t.Errorf("Describe() sent %d descriptors, want %d", len(descs), expectedCount)
	}
}

func TestAirGradientCollector_Collect_Success(t *testing.T) {
	expectedMeasures := &airgradient.Measures{
		Pm01:            10.5,
		Pm02:            15.3,
		Pm10:            20.1,
		Pm01Standard:    11.0,
		Pm02Standard:    16.0,
		Pm10Standard:    21.0,
		Pm003Count:      1000.0,
		Pm005Count:      2000.0,
		Pm01Count:       3000.0,
		Pm02Count:       4000.0,
		Pm50Count:       5000.0,
		Pm10Count:       6000.0,
		Pm02Compensated: 15.5,
		Atmp:            22.5,
		AtmpCompensated: 22.7,
		Rhum:            45.0,
		RhumCompensated: 45.2,
		Rco2:            420,
		TvocIndex:       100,
		TvocRaw:         200,
		NoxIndex:        50,
		NoxRaw:          150,
		Boot:            12345,
		BootCount:       5,
		Wifi:            -65,
		SerialNo:        "TEST123",
		Firmware:        "1.0.0",
		Model:           "TEST",
	}

	expectedConfig := &airgradient.Config{
		Country:           "US",
		PmStandard:        "EPA",
		LedBarMode:        "auto",
		TemperatureUnit:   "C",
		Model:             "TEST",
		DisplayBrightness: 80,
		LedBarBrightness:  50,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/measures/current" {
			json.NewEncoder(w).Encode(expectedMeasures)
		} else if r.URL.Path == "/config" {
			json.NewEncoder(w).Encode(expectedConfig)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	host := server.URL[7:] // Remove "http://" prefix
	collector := NewAirGradientCollector(logger, []string{host})

	ch := make(chan prometheus.Metric, 100)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	// Should have up=1 + all measures metrics + all config metrics
	expectedMetricCount := 1 + 25 + 3 // up + measures + config
	if len(metrics) < expectedMetricCount {
		t.Errorf("Collect() sent %d metrics, want at least %d", len(metrics), expectedMetricCount)
	}

	// Verify up metric
	var upMetric prometheus.Metric
	for _, m := range metrics {
		desc := m.Desc()
		if desc.String() == collector.up.String() {
			upMetric = m
			break
		}
	}

	if upMetric == nil {
		t.Fatal("up metric not found")
	}

	var dtoMetric dto.Metric
	err := upMetric.Write(&dtoMetric)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if dtoMetric.Gauge == nil || *dtoMetric.Gauge.Value != 1.0 {
		t.Errorf("up metric value = %v, want 1.0", dtoMetric.Gauge.Value)
	}
}

func TestAirGradientCollector_Collect_MeasuresError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/measures/current" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	host := server.URL[7:] // Remove "http://" prefix
	collector := NewAirGradientCollector(logger, []string{host})

	ch := make(chan prometheus.Metric, 100)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	// Should only have up=0 metric
	if len(metrics) != 1 {
		t.Errorf("Collect() sent %d metrics, want 1", len(metrics))
	}

	var dtoMetric dto.Metric
	err := metrics[0].Write(&dtoMetric)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if dtoMetric.Gauge == nil || *dtoMetric.Gauge.Value != 0.0 {
		t.Errorf("up metric value = %v, want 0.0", dtoMetric.Gauge.Value)
	}
}

func TestAirGradientCollector_Collect_ConfigError(t *testing.T) {
	expectedMeasures := &airgradient.Measures{
		Pm01:     10.5,
		Pm02:     15.3,
		SerialNo: "TEST123",
		Firmware: "1.0.0",
		Model:    "TEST",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/measures/current" {
			json.NewEncoder(w).Encode(expectedMeasures)
		} else if r.URL.Path == "/config" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	host := server.URL[7:] // Remove "http://" prefix
	collector := NewAirGradientCollector(logger, []string{host})

	ch := make(chan prometheus.Metric, 100)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	// Should have up=1 + all measures metrics (but no config metrics)
	expectedMetricCount := 1 + 25 // up + measures
	if len(metrics) < expectedMetricCount {
		t.Errorf("Collect() sent %d metrics, want at least %d", len(metrics), expectedMetricCount)
	}
}

func TestAirGradientCollector_Collect_MultipleHosts(t *testing.T) {
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/measures/current" {
			json.NewEncoder(w).Encode(&airgradient.Measures{
				Pm01:     10.5,
				SerialNo: "TEST1",
				Firmware: "1.0.0",
				Model:    "TEST",
			})
		} else if r.URL.Path == "/config" {
			json.NewEncoder(w).Encode(&airgradient.Config{
				Model: "TEST",
			})
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/measures/current" {
			json.NewEncoder(w).Encode(&airgradient.Measures{
				Pm01:     20.5,
				SerialNo: "TEST2",
				Firmware: "1.0.0",
				Model:    "TEST",
			})
		} else if r.URL.Path == "/config" {
			json.NewEncoder(w).Encode(&airgradient.Config{
				Model: "TEST",
			})
		}
	}))
	defer server2.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	host1 := server1.URL[7:] // Remove "http://" prefix
	host2 := server2.URL[7:] // Remove "http://" prefix
	collector := NewAirGradientCollector(logger, []string{host1, host2})

	ch := make(chan prometheus.Metric, 200)
	collector.Collect(ch)
	close(ch)

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	// Should have 2 * (1 + 25 + 3) metrics (up + measures + config for each host)
	expectedMetricCount := 2 * (1 + 25 + 3)
	if len(metrics) < expectedMetricCount {
		t.Errorf("Collect() sent %d metrics, want at least %d", len(metrics), expectedMetricCount)
	}
}
