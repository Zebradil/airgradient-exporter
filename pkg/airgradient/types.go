package airgradient

type Measures struct {
	Pm01            float64 `json:"pm01"`
	Pm02            float64 `json:"pm02"`
	Pm10            float64 `json:"pm10"`
	Pm01Standard    float64 `json:"pm01Standard"`
	Pm02Standard    float64 `json:"pm02Standard"`
	Pm10Standard    float64 `json:"pm10Standard"`
	Pm003Count      float64 `json:"pm003Count"`
	Pm005Count      float64 `json:"pm005Count"`
	Pm01Count       float64 `json:"pm01Count"`
	Pm02Count       float64 `json:"pm02Count"`
	Pm50Count       float64 `json:"pm50Count"`
	Pm10Count       float64 `json:"pm10Count"`
	Pm02Compensated float64 `json:"pm02Compensated"`
	Atmp            float64 `json:"atmp"`
	AtmpCompensated float64 `json:"atmpCompensated"`
	Rhum            float64 `json:"rhum"`
	RhumCompensated float64 `json:"rhumCompensated"`
	Rco2            float64 `json:"rco2"`
	TvocIndex       float64 `json:"tvocIndex"`
	TvocRaw         float64 `json:"tvocRaw"`
	NoxIndex        float64 `json:"noxIndex"`
	NoxRaw          float64 `json:"noxRaw"`
	Boot            int     `json:"boot"`
	BootCount       int     `json:"bootCount"`
	Wifi            int     `json:"wifi"`
	LedMode         string  `json:"ledMode"`
	SerialNo        string  `json:"serialno"`
	Firmware        string  `json:"firmware"`
	Model           string  `json:"model"`
}

type Config struct {
	Country                         string      `json:"country"`
	PmStandard                      string      `json:"pmStandard"`
	LedBarMode                      string      `json:"ledBarMode"`
	AbcDays                         int         `json:"abcDays"`
	TvocLearningOffset              int         `json:"tvocLearningOffset"`
	NoxLearningOffset               int         `json:"noxLearningOffset"`
	MqttBrokerURL                   string      `json:"mqttBrokerUrl"`
	HTTPDomain                      string      `json:"httpDomain"`
	TemperatureUnit                 string      `json:"temperatureUnit"`
	DisableCloudConnection          bool        `json:"disableCloudConnection"`
	ConfigurationControl            string      `json:"configurationControl"`
	PostDataToAirGradient           bool        `json:"postDataToAirGradient"`
	LedBarBrightness                int         `json:"ledBarBrightness"`
	DisplayBrightness               int         `json:"displayBrightness"`
	OfflineMode                     bool        `json:"offlineMode"`
	MonitorDisplayCompensatedValues bool        `json:"monitorDisplayCompensatedValues"`
	Model                           string      `json:"model"`
	Corrections                     Corrections `json:"corrections"`
}

type Corrections struct {
	Pm02 ParamCorrection `json:"pm02"`
}

type ParamCorrection struct {
	CorrectionAlgorithm string `json:"correctionAlgorithm"`
	Slr                 any    `json:"slr"` // value is null in example
}
