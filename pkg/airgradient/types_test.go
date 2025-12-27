package airgradient

import (
	"encoding/json"
	"testing"
)

func TestMeasures_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"atmp": 22.5,
		"atmpCompensated": 22.7,
		"boot": 12345,
		"bootCount": 5,
		"firmware": "1.0.0",
		"ledMode": "auto",
		"model": "TEST",
		"noxIndex": 50,
		"noxRaw": 150,
		"pm003Count": 1000.0,
		"pm005Count": 2000.0,
		"pm01": 10.5,
		"pm01Count": 3000.0,
		"pm01Standard": 11.0,
		"pm02": 15.3,
		"pm02Compensated": 15.5,
		"pm02Count": 4000.0,
		"pm02Standard": 16.0,
		"pm10": 20.1,
		"pm10Count": 6000.0,
		"pm10Standard": 21.0,
		"pm50Count": 5000.0,
		"rco2": 420,
		"rhum": 45.0,
		"rhumCompensated": 45.2,
		"serialno": "TEST123",
		"tvocIndex": 100,
		"tvocRaw": 200,
		"wifi": -65
	}`

	var measures Measures
	err := json.Unmarshal([]byte(jsonData), &measures)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if measures.Pm01 != 10.5 {
		t.Errorf("Pm01 = %v, want 10.5", measures.Pm01)
	}
	if measures.Pm02 != 15.3 {
		t.Errorf("Pm02 = %v, want 15.3", measures.Pm02)
	}
	if measures.SerialNo != "TEST123" {
		t.Errorf("SerialNo = %v, want TEST123", measures.SerialNo)
	}
	if measures.Boot != 12345 {
		t.Errorf("Boot = %v, want 12345", measures.Boot)
	}
	if measures.Wifi != -65 {
		t.Errorf("Wifi = %v, want -65", measures.Wifi)
	}
}

func TestConfig_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"abcDays": 7,
		"configurationControl": "local",
		"country": "US",
		"disableCloudConnection": false,
		"displayBrightness": 80,
		"httpDomain": "example.com",
		"ledBarBrightness": 50,
		"ledBarMode": "auto",
		"model": "TEST",
		"monitorDisplayCompensatedValues": true,
		"mqttBrokerUrl": "mqtt://example.com",
		"noxLearningOffset": 200,
		"offlineMode": false,
		"pmStandard": "EPA",
		"postDataToAirGradient": true,
		"temperatureUnit": "C",
		"tvocLearningOffset": 100,
		"corrections": {
			"pm02": {
				"correctionAlgorithm": "slr",
				"slr": null
			}
		}
	}`

	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if config.Country != "US" {
		t.Errorf("Country = %v, want US", config.Country)
	}
	if config.PmStandard != "EPA" {
		t.Errorf("PmStandard = %v, want EPA", config.PmStandard)
	}
	if config.Model != "TEST" {
		t.Errorf("Model = %v, want TEST", config.Model)
	}
	if config.DisplayBrightness != 80 {
		t.Errorf("DisplayBrightness = %v, want 80", config.DisplayBrightness)
	}
	if config.Corrections.Pm02.CorrectionAlgorithm != "slr" {
		t.Errorf("CorrectionAlgorithm = %v, want slr", config.Corrections.Pm02.CorrectionAlgorithm)
	}
}

func TestConfig_UnmarshalJSON_WithSLRValue(t *testing.T) {
	jsonData := `{
		"country": "US",
		"model": "TEST",
		"pmStandard": "EPA",
		"corrections": {
			"pm02": {
				"correctionAlgorithm": "slr",
				"slr": 1.5
			}
		}
	}`

	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if config.Corrections.Pm02.Slr != 1.5 {
		t.Errorf("Slr = %v, want 1.5", config.Corrections.Pm02.Slr)
	}
}

func TestMeasures_UnmarshalJSON_Minimal(t *testing.T) {
	jsonData := `{
		"firmware": "",
		"model": "",
		"pm01": 0,
		"pm02": 0,
		"pm10": 0,
		"serialno": ""
	}`

	var measures Measures
	err := json.Unmarshal([]byte(jsonData), &measures)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if measures.Pm01 != 0 {
		t.Errorf("Pm01 = %v, want 0", measures.Pm01)
	}
}
