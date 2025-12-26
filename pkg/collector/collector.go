package collector

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/zebradil/airgradient-exporter/pkg/airgradient"
)

const (
	namespace = "airgradient"
)

type AirGradientCollector struct {
	client *airgradient.Client
	hosts  []string
	logger *slog.Logger

	up              *prometheus.Desc
	pm01            *prometheus.Desc
	pm02            *prometheus.Desc
	pm10            *prometheus.Desc
	pm01Standard    *prometheus.Desc
	pm02Standard    *prometheus.Desc
	pm10Standard    *prometheus.Desc
	pm003Count      *prometheus.Desc
	pm005Count      *prometheus.Desc
	pm01Count       *prometheus.Desc
	pm02Count       *prometheus.Desc
	pm50Count       *prometheus.Desc
	pm10Count       *prometheus.Desc
	pm02Compensated *prometheus.Desc
	atmp            *prometheus.Desc
	atmpCompensated *prometheus.Desc
	rhum            *prometheus.Desc
	rhumCompensated *prometheus.Desc
	rco2            *prometheus.Desc
	tvocIndex       *prometheus.Desc
	tvocRaw         *prometheus.Desc
	noxIndex        *prometheus.Desc
	noxRaw          *prometheus.Desc
	boot            *prometheus.Desc
	bootCount       *prometheus.Desc
	wifi            *prometheus.Desc

	// Config metrics
	displayBrightness *prometheus.Desc
	ledBarBrightness  *prometheus.Desc
	configInfo        *prometheus.Desc
}

func NewAirGradientCollector(logger *slog.Logger, hosts []string) *AirGradientCollector {
	return &AirGradientCollector{
		client: airgradient.NewClient(),
		hosts:  hosts,
		logger: logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Was the last scrape of the AirGradient monitor successful.",
			[]string{"host"}, nil,
		),
		pm01: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm01"),
			"PM1.0 (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm02: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm02"),
			"PM2.5 (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm10: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm10"),
			"PM10 (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm01Standard: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm01_standard"),
			"PM1.0 Standard (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm02Standard: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm02_standard"),
			"PM2.5 Standard (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm10Standard: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm10_standard"),
			"PM10 Standard (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm003Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm003_count"),
			"PM0.3 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm005Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm005_count"),
			"PM0.5 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm01Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm01_count"),
			"PM1.0 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm02Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm02_count"),
			"PM2.5 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm50Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm50_count"),
			"PM5.0 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm10Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm10_count"),
			"PM10 Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		pm02Compensated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "pm02_compensated"),
			"PM2.5 Compensated (ug/m3)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		atmp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "atmp"),
			"Ambient Temperature (Celsius)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		atmpCompensated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "atmp_compensated"),
			"Ambient Temperature Compensated (Celsius)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		rhum: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "rhum"),
			"Relative Humidity (%)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		rhumCompensated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "rhum_compensated"),
			"Relative Humidity Compensated (%)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		rco2: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "rco2"),
			"CO2 (ppm)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		tvocIndex: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "tvoc_index"),
			"TVOC Index",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		tvocRaw: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "tvoc_raw"),
			"TVOC Raw",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		noxIndex: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nox_index"),
			"NOx Index",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		noxRaw: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nox_raw"),
			"NOx Raw",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		boot: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_uptime_raw"),
			"Boot uptime raw value",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		bootCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_count"),
			"Boot Count",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		wifi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "wifi_signal_dbm"),
			"WiFi Signal (dBm)",
			[]string{"host", "serialno", "firmware"}, nil,
		),
		displayBrightness: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "config_display_brightness"),
			"Display Brightness",
			[]string{"host", "model"}, nil,
		),
		ledBarBrightness: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "config_led_bar_brightness"),
			"LED Bar Brightness",
			[]string{"host", "model"}, nil,
		),
		configInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "config_info"),
			"Configuration Info",
			[]string{"host", "model", "country", "pm_standard", "led_bar_mode", "temperature_unit"}, nil,
		),
	}
}

func (c *AirGradientCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.pm01
	ch <- c.pm02
	ch <- c.pm10
	ch <- c.pm01Standard
	ch <- c.pm02Standard
	ch <- c.pm10Standard
	ch <- c.pm003Count
	ch <- c.pm005Count
	ch <- c.pm01Count
	ch <- c.pm02Count
	ch <- c.pm50Count
	ch <- c.pm10Count
	ch <- c.pm02Compensated
	ch <- c.atmp
	ch <- c.atmpCompensated
	ch <- c.rhum
	ch <- c.rhumCompensated
	ch <- c.rco2
	ch <- c.tvocIndex
	ch <- c.tvocRaw
	ch <- c.noxIndex
	ch <- c.noxRaw
	ch <- c.boot
	ch <- c.bootCount
	ch <- c.wifi
	ch <- c.displayBrightness
	ch <- c.ledBarBrightness
	ch <- c.configInfo
}

func (c *AirGradientCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(len(c.hosts))

	for _, host := range c.hosts {
		go func(host string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c.scrape(ctx, host, ch)
		}(host)
	}

	wg.Wait()
}

func (c *AirGradientCollector) scrape(ctx context.Context, host string, ch chan<- prometheus.Metric) {
	// We need a separate context for each scrape or handle cancellation?
	// Actually, context should be passed in. `Collect` doesn't take context.
	// So we create one.

	// Fetch Measures
	measures, err := c.client.GetMeasures(ctx, host)
	if err != nil {
		c.logger.Error("Failed to scrape measures", "host", host, "error", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0, host)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1, host)

	labels := []string{host, measures.SerialNo, measures.Firmware}

	ch <- prometheus.MustNewConstMetric(c.pm01, prometheus.GaugeValue, measures.Pm01, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm02, prometheus.GaugeValue, measures.Pm02, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm10, prometheus.GaugeValue, measures.Pm10, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm01Standard, prometheus.GaugeValue, measures.Pm01Standard, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm02Standard, prometheus.GaugeValue, measures.Pm02Standard, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm10Standard, prometheus.GaugeValue, measures.Pm10Standard, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm003Count, prometheus.GaugeValue, measures.Pm003Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm005Count, prometheus.GaugeValue, measures.Pm005Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm01Count, prometheus.GaugeValue, measures.Pm01Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm02Count, prometheus.GaugeValue, measures.Pm02Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm50Count, prometheus.GaugeValue, measures.Pm50Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm10Count, prometheus.GaugeValue, measures.Pm10Count, labels...)
	ch <- prometheus.MustNewConstMetric(c.pm02Compensated, prometheus.GaugeValue, measures.Pm02Compensated, labels...)
	ch <- prometheus.MustNewConstMetric(c.atmp, prometheus.GaugeValue, measures.Atmp, labels...)
	ch <- prometheus.MustNewConstMetric(c.atmpCompensated, prometheus.GaugeValue, measures.AtmpCompensated, labels...)
	ch <- prometheus.MustNewConstMetric(c.rhum, prometheus.GaugeValue, measures.Rhum, labels...)
	ch <- prometheus.MustNewConstMetric(c.rhumCompensated, prometheus.GaugeValue, measures.RhumCompensated, labels...)
	ch <- prometheus.MustNewConstMetric(c.rco2, prometheus.GaugeValue, measures.Rco2, labels...)
	ch <- prometheus.MustNewConstMetric(c.tvocIndex, prometheus.GaugeValue, measures.TvocIndex, labels...)
	ch <- prometheus.MustNewConstMetric(c.tvocRaw, prometheus.GaugeValue, measures.TvocRaw, labels...)
	ch <- prometheus.MustNewConstMetric(c.noxIndex, prometheus.GaugeValue, measures.NoxIndex, labels...)
	ch <- prometheus.MustNewConstMetric(c.noxRaw, prometheus.GaugeValue, measures.NoxRaw, labels...)
	ch <- prometheus.MustNewConstMetric(c.boot, prometheus.GaugeValue, float64(measures.Boot), labels...)
	ch <- prometheus.MustNewConstMetric(c.bootCount, prometheus.GaugeValue, float64(measures.BootCount), labels...)
	ch <- prometheus.MustNewConstMetric(c.wifi, prometheus.GaugeValue, float64(measures.Wifi), labels...)

	// Fetch Config
	config, err := c.client.GetConfig(ctx, host)
	if err != nil {
		c.logger.Error("Failed to scrape config", "host", host, "error", err)
		return
	}

	configLabels := []string{host, config.Model}
	ch <- prometheus.MustNewConstMetric(c.displayBrightness, prometheus.GaugeValue, float64(config.DisplayBrightness), configLabels...)
	ch <- prometheus.MustNewConstMetric(c.ledBarBrightness, prometheus.GaugeValue, float64(config.LedBarBrightness), configLabels...)

	infoLabels := []string{host, config.Model, config.Country, config.PmStandard, config.LedBarMode, config.TemperatureUnit}
	ch <- prometheus.MustNewConstMetric(c.configInfo, prometheus.GaugeValue, 1, infoLabels...)
}
