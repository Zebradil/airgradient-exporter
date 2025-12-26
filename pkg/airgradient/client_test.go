package airgradient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_GetMeasures_Success(t *testing.T) {
	expectedMeasures := &Measures{
		Pm01:      10.5,
		Pm02:      15.3,
		Pm10:      20.1,
		Atmp:      22.5,
		Rhum:      45.0,
		Rco2:      420,
		SerialNo:  "TEST123",
		Firmware:  "1.0.0",
		Model:     "TEST",
		Boot:      12345,
		BootCount: 5,
		Wifi:      -65,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/measures/current" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedMeasures)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	measures, err := client.GetMeasures(context.Background(), host)
	if err != nil {
		t.Fatalf("GetMeasures() error = %v", err)
	}

	if measures == nil {
		t.Fatal("GetMeasures() returned nil")
	}

	if measures.Pm01 != expectedMeasures.Pm01 {
		t.Errorf("Pm01 = %v, want %v", measures.Pm01, expectedMeasures.Pm01)
	}
	if measures.SerialNo != expectedMeasures.SerialNo {
		t.Errorf("SerialNo = %v, want %v", measures.SerialNo, expectedMeasures.SerialNo)
	}
}

func TestClient_GetMeasures_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	measures, err := client.GetMeasures(context.Background(), host)
	if err == nil {
		t.Fatal("GetMeasures() expected error, got nil")
	}
	if measures != nil {
		t.Fatal("GetMeasures() expected nil measures on error")
	}
}

func TestClient_GetMeasures_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	measures, err := client.GetMeasures(context.Background(), host)
	if err == nil {
		t.Fatal("GetMeasures() expected error, got nil")
	}
	if measures != nil {
		t.Fatal("GetMeasures() expected nil measures on error")
	}
}

func TestClient_GetMeasures_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	measures, err := client.GetMeasures(ctx, host)
	if err == nil {
		t.Fatal("GetMeasures() expected error due to timeout, got nil")
	}
	if measures != nil {
		t.Fatal("GetMeasures() expected nil measures on error")
	}
}

func TestClient_GetMeasures_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	measures, err := client.GetMeasures(ctx, host)
	if err == nil {
		t.Fatal("GetMeasures() expected error due to cancellation, got nil")
	}
	if measures != nil {
		t.Fatal("GetMeasures() expected nil measures on error")
	}
}

func TestClient_GetConfig_Success(t *testing.T) {
	expectedConfig := &Config{
		Country:                "US",
		PmStandard:             "EPA",
		LedBarMode:             "auto",
		TemperatureUnit:        "C",
		Model:                  "TEST",
		DisplayBrightness:      80,
		LedBarBrightness:       50,
		AbcDays:                7,
		DisableCloudConnection: false,
		PostDataToAirGradient:  true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedConfig)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	config, err := client.GetConfig(context.Background(), host)
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}

	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}

	if config.Country != expectedConfig.Country {
		t.Errorf("Country = %v, want %v", config.Country, expectedConfig.Country)
	}
	if config.Model != expectedConfig.Model {
		t.Errorf("Model = %v, want %v", config.Model, expectedConfig.Model)
	}
	if config.DisplayBrightness != expectedConfig.DisplayBrightness {
		t.Errorf("DisplayBrightness = %v, want %v", config.DisplayBrightness, expectedConfig.DisplayBrightness)
	}
}

func TestClient_GetConfig_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	config, err := client.GetConfig(context.Background(), host)
	if err == nil {
		t.Fatal("GetConfig() expected error, got nil")
	}
	if config != nil {
		t.Fatal("GetConfig() expected nil config on error")
	}
}

func TestClient_GetConfig_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	config, err := client.GetConfig(context.Background(), host)
	if err == nil {
		t.Fatal("GetConfig() expected error, got nil")
	}
	if config != nil {
		t.Fatal("GetConfig() expected nil config on error")
	}
}

func TestClient_GetConfig_WithCorrections(t *testing.T) {
	expectedConfig := &Config{
		Country:    "US",
		PmStandard: "EPA",
		Model:      "TEST",
		Corrections: Corrections{
			Pm02: ParamCorrection{
				CorrectionAlgorithm: "slr",
				Slr:                 nil,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedConfig)
	}))
	defer server.Close()

	client := NewClient()
	host := server.URL[7:] // Remove "http://" prefix

	config, err := client.GetConfig(context.Background(), host)
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}

	if config.Corrections.Pm02.CorrectionAlgorithm != "slr" {
		t.Errorf("CorrectionAlgorithm = %v, want slr", config.Corrections.Pm02.CorrectionAlgorithm)
	}
}
