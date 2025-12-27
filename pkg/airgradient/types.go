package airgradient

type Measures struct {
	Atmp            float64 `json:"atmp"`
	AtmpCompensated float64 `json:"atmpCompensated"`
	Boot            int     `json:"boot"`
	BootCount       int     `json:"bootCount"`
	Firmware        string  `json:"firmware"`
	LedMode         string  `json:"ledMode"`
	Model           string  `json:"model"`
	NoxIndex        float64 `json:"noxIndex"`
	NoxRaw          float64 `json:"noxRaw"`
	Pm003Count      float64 `json:"pm003Count"`
	Pm005Count      float64 `json:"pm005Count"`
	Pm01            float64 `json:"pm01"`
	Pm01Count       float64 `json:"pm01Count"`
	Pm01Standard    float64 `json:"pm01Standard"`
	Pm02            float64 `json:"pm02"`
	Pm02Compensated float64 `json:"pm02Compensated"`
	Pm02Count       float64 `json:"pm02Count"`
	Pm02Standard    float64 `json:"pm02Standard"`
	Pm10            float64 `json:"pm10"`
	Pm10Count       float64 `json:"pm10Count"`
	Pm10Standard    float64 `json:"pm10Standard"`
	Pm50Count       float64 `json:"pm50Count"`
	Rco2            float64 `json:"rco2"`
	Rhum            float64 `json:"rhum"`
	RhumCompensated float64 `json:"rhumCompensated"`
	SerialNo        string  `json:"serialno"`
	TvocIndex       float64 `json:"tvocIndex"`
	TvocRaw         float64 `json:"tvocRaw"`
	Wifi            int     `json:"wifi"`
}

type Config struct {
	AbcDays                         int         `json:"abcDays"`
	ConfigurationControl            string      `json:"configurationControl"`
	Corrections                     Corrections `json:"corrections"`
	Country                         string      `json:"country"`
	DisableCloudConnection          bool        `json:"disableCloudConnection"`
	DisplayBrightness               int         `json:"displayBrightness"`
	HTTPDomain                      string      `json:"httpDomain"`
	LedBarBrightness                int         `json:"ledBarBrightness"`
	LedBarMode                      string      `json:"ledBarMode"`
	Model                           string      `json:"model"`
	MonitorDisplayCompensatedValues bool        `json:"monitorDisplayCompensatedValues"`
	MqttBrokerURL                   string      `json:"mqttBrokerUrl"`
	NoxLearningOffset               int         `json:"noxLearningOffset"`
	OfflineMode                     bool        `json:"offlineMode"`
	PmStandard                      string      `json:"pmStandard"`
	PostDataToAirGradient           bool        `json:"postDataToAirGradient"`
	TemperatureUnit                 string      `json:"temperatureUnit"`
	TvocLearningOffset              int         `json:"tvocLearningOffset"`
}

type Corrections struct {
	Pm02 ParamCorrection `json:"pm02"`
}

type ParamCorrection struct {
	CorrectionAlgorithm string `json:"correctionAlgorithm"`
	Slr                 any    `json:"slr"` // value is null in example
}
