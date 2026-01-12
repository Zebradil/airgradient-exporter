package main

import (
	"fmt"
	"sort"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

// ============================================================================
// Definitions
// ============================================================================

type metricDefinition struct {
	metric      string
	title       string
	unit        string
	description string
	thresholds  thresholdDefinition
}

type thresholdDefinition map[float64]string

// Grafana automatically adds a "80: red" threshold step if none are defined,
// so we define an explicit "no thresholds" definition for metrics without thresholds.
var noThresholds = thresholdDefinition{
	// 0:     "green",
	// 10000: "green",
}

var metricTemperature = metricDefinition{
	metric:      "airgradient_atmp",
	title:       "Temperature (°C)",
	unit:        "celsius",
	description: "Ambient temperature near the sensor.",
	thresholds: thresholdDefinition{
		0:  "semi-dark-blue",
		18: "light-blue",
		21: "green",
		25: "yellow",
		28: "red",
	},
}

var metricHumidity = metricDefinition{
	metric:      "airgradient_rhum",
	title:       "Relative Humidity (%)",
	unit:        "percent",
	description: "Ambient relative humidity near the sensor.",
	thresholds: thresholdDefinition{
		0:  "orange",
		25: "yellow",
		30: "green",
		50: "yellow",
		60: "orange",
	},
}

var metricCO2 = metricDefinition{
	metric:      "airgradient_rco2",
	title:       "CO2 Level (ppm)",
	unit:        "ppm",
	description: "Carbon dioxide concentration in the air. PPM (parts per million) reflects volumetric gas concentration. Standard unit for tracking CO2 buildup from human activity.",
	thresholds: thresholdDefinition{
		0:     "dark-green",
		801:   "green",
		1000:  "yellow",
		1500:  "orange",
		2000:  "red",
		3000:  "purple",
		10000: "#000000",
	},
}

var metricPM1 = metricDefinition{
	metric:      "airgradient_pm01",
	title:       "PM1 μg/m3",
	unit:        "conμgm3",
	description: "Total weight of particles with a diameter of 1.0 micrometers or smaller, in μg / m3",
	thresholds:  noThresholds,
}

var metricPM25 = metricDefinition{
	metric:      "airgradient_pm02",
	title:       "PM2.5 μg/m3",
	unit:        "conμgm3",
	description: "Total weight of particles with a diameter of 2.5 micrometers or smaller, in μg / m3",
	thresholds:  noThresholds,
}

var metricPM10 = metricDefinition{
	metric:      "airgradient_pm10",
	title:       "PM10 μg/m3",
	unit:        "conμgm3",
	description: "Total weight of particles with a diameter of 10.0 micrometers or smaller, in μg / m3",
	thresholds:  noThresholds,
}

var metricPM03 = metricDefinition{
	metric:      "airgradient_pm003_count",
	title:       "PM0.3 particles/dL",
	unit:        "particles/100ml",
	description: "Total number of particles with a diameter of 0.3 micrometers or smaller per 100ml air",
	thresholds:  noThresholds,
}

var metricNOx = metricDefinition{
	metric:      "airgradient_nox_index",
	title:       "NOx Index",
	unit:        "NOx",
	description: "Unitless index (1 - 500) based on deviation from 24h baseline. 1 = baseline, >1 = more oxidizing gases present.",
	thresholds: thresholdDefinition{
		0:   "green",
		1.1: "red",
	},
}

var metricTVOC = metricDefinition{
	metric:      "airgradient_tvoc_index",
	title:       "TVOC Index",
	unit:        "TVOC",
	description: "Unitless index (1 - 500) based on deviation from 24h baseline. 100 = baseline, >100 = more VOCs present compared to the 24-hour baseline, <100 = fewer VOCs compared to the 24-hour baseline. Index is used because MOx sensors respond broadly to many VOCs, making exact ppb estimates unreliable. Better for detecting relative changes in air quality. The VOC Index adapts its gain based on environment: in clean air it becomes more sensitive to small changes, and in polluted air, it reduces sensitivity to avoid saturation, ensuring relative changes remain meaningful.",
	thresholds: thresholdDefinition{
		0:   "green",
		100: "yellow",
		101: "red",
	},
}

// ============================================================================
// Main
// ============================================================================

func buildDashboard() (dashboard.Dashboard, error) {
	builder := dashboard.NewDashboardBuilder("Airgradient Native").
		Uid("d11cb5b1-cb9b-47a9-98a6-94cc40fbc0a6").
		Description("Dashboard for the Airgradient One native prometheus metrics exporter. Compatible with the 3.0.6 firmware.").
		Editable().
		LiveNow(true).
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Timezone("").
		Time("now-12h", "now").
		Refresh("").
		// Variables
		WithVariable(buildSerialNoVariable()).
		// Current Status Row
		WithRow(buildCurrentStatusRow()).
		// Stat Panels
		WithPanel(buildTemperatureStatPanel()).
		WithPanel(buildHumidityStatPanel()).
		WithPanel(buildCO2StatPanel()).
		WithPanel(buildPM1StatPanel()).
		WithPanel(buildPM25StatPanel()).
		WithPanel(buildPM10StatPanel()).
		WithPanel(buildPM03StatPanel()).
		WithPanel(buildNOxStatPanel()).
		WithPanel(buildTVOCStatPanel()).
		// Temperature & Humidity Row
		WithRow(buildTempHumidityRow()).
		WithPanel(buildTemperatureTimeseriesPanel()).
		WithPanel(buildHumidityTimeseriesPanel()).
		// CO2 Row
		WithRow(buildCO2Row()).
		WithPanel(buildCO2TimeseriesPanel()).
		// Particle Pollution Row
		WithRow(buildParticlePollutionRow()).
		WithPanel(buildPM03TimeseriesPanel()).
		WithPanel(buildPM1TimeseriesPanel()).
		WithPanel(buildPM25TimeseriesPanel()).
		WithPanel(buildPM10TimeseriesPanel()).
		// TVOC / NOx Row
		WithRow(buildTVOCNOxRow()).
		WithPanel(buildTVOCTimeseriesPanel()).
		WithPanel(buildNOxTimeseriesPanel())

	return builder.Build()
}

// ============================================================================
// Variable
// ============================================================================

func buildSerialNoVariable() *dashboard.QueryVariableBuilder {
	return dashboard.NewQueryVariableBuilder("serialno").
		Datasource(dashboard.DataSourceRef{
			Type: cog.ToPtr("prometheus"),
		}).
		Query(dashboard.StringOrMap{
			String: cog.ToPtr("label_values(airgradient_config_info,serialno)"),
		}).
		AllValue(".*").
		IncludeAll(true).
		Multi(true).
		Refresh(dashboard.VariableRefreshOnDashboardLoad).
		Sort(dashboard.VariableSortAlphabeticalAsc)
}

// ============================================================================
// Row Panels
// ============================================================================

func buildCurrentStatusRow() *dashboard.RowBuilder {
	return dashboard.NewRowBuilder("Current Status (serial $serialno)").
		GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 0}).
		Repeat("serialno")
}

func buildTempHumidityRow() *dashboard.RowBuilder {
	return dashboard.NewRowBuilder("Temperature & Humidity").
		GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 9})
}

func buildCO2Row() *dashboard.RowBuilder {
	return dashboard.NewRowBuilder("CO2 PPM").
		GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 23})
}

func buildParticlePollutionRow() *dashboard.RowBuilder {
	return dashboard.NewRowBuilder("Particle Pollution").
		GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 37})
}

func buildTVOCNOxRow() *dashboard.RowBuilder {
	return dashboard.NewRowBuilder("TVOC / NOx").
		GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 66})
}

// ============================================================================
// Stat Panels - Current Status
// ============================================================================

func buildStatPanel(metric metricDefinition, pos dashboard.GridPos) *stat.PanelBuilder {
	builder := stat.NewPanelBuilder().
		Title(metric.title).
		Description(metric.description).
		GridPos(pos).
		Datasource(prometheusDatasourceRef()).
		WithTarget(
			prometheusQuery(fmt.Sprintf(`sum by(serialno) (%s{serialno=~"$serialno"})`, metric.metric), "{{sensor}}").
				Instant(),
		).
		Thresholds(buildThresholds(metric.thresholds)).
		ColorMode(common.BigValueColorModeValue).
		ColorScheme(
			dashboard.NewFieldColorBuilder().
				Mode(dashboard.FieldColorModeIdPaletteClassic).
				SeriesBy(dashboard.FieldColorSeriesByModeLast),
		).
		GraphMode(common.BigValueGraphModeNone).
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}).
				Values(true),
		)
	if len(metric.thresholds) > 0 {
		builder = builder.
			ColorMode(common.BigValueColorModeBackground).
			ColorScheme(
				dashboard.NewFieldColorBuilder().
					Mode(dashboard.FieldColorModeIdThresholds).
					SeriesBy(dashboard.FieldColorSeriesByModeLast),
			)
	}
	return builder
}

func buildTemperatureStatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricTemperature, dashboard.GridPos{H: 4, W: 4, X: 0, Y: 1})
}

func buildHumidityStatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricHumidity, dashboard.GridPos{H: 4, W: 4, X: 4, Y: 1})
}

func buildCO2StatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricCO2, dashboard.GridPos{H: 4, W: 4, X: 8, Y: 1})
}

func buildPM1StatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricPM1, dashboard.GridPos{H: 4, W: 4, X: 12, Y: 1})
}

func buildPM25StatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricPM25, dashboard.GridPos{H: 4, W: 4, X: 16, Y: 1})
}

func buildPM10StatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricPM10, dashboard.GridPos{H: 4, W: 4, X: 20, Y: 1})
}

func buildPM03StatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricPM03, dashboard.GridPos{H: 4, W: 4, X: 0, Y: 5})
}

func buildNOxStatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricNOx, dashboard.GridPos{H: 4, W: 4, X: 4, Y: 5})
}

func buildTVOCStatPanel() *stat.PanelBuilder {
	return buildStatPanel(metricTVOC, dashboard.GridPos{H: 4, W: 4, X: 8, Y: 5})
}

// ============================================================================
// Timeseries Panels - Temperature & Humidity
// ============================================================================

func buildTimeseriesPanel(metric metricDefinition, pos dashboard.GridPos) *timeseries.PanelBuilder {
	builder := timeseries.NewPanelBuilder().
		Title(metric.title).
		Description(metric.description).
		Unit(metric.unit).
		GridPos(pos).
		Datasource(prometheusDatasourceRef()).
		WithTarget(
			prometheusQuery(fmt.Sprintf(`sum by(serialno) (%s{serialno=~"$serialno"})`, metric.metric), "{{sensor}}").
				Range(),
		).
		Thresholds(buildThresholds(metric.thresholds)).
		LineWidth(2).
		FillOpacity(15).
		LineInterpolation(common.LineInterpolationSmooth).
		GradientMode(common.GraphGradientModeHue).
		ShowPoints(common.VisibilityModeAuto).
		SpanNulls(common.BoolOrFloat64{Bool: cog.ToPtr(false)}).
		AxisBorderShow(true).
		Legend(
			common.NewVizLegendOptionsBuilder().
				DisplayMode(common.LegendDisplayModeTable).
				Placement(common.LegendPlacementBottom).
				ShowLegend(true).
				Calcs([]string{"lastNotNull", "min", "max", "mean", "variance", "stdDev"}),
		).
		Tooltip(
			common.NewVizTooltipOptionsBuilder().
				Mode(common.TooltipDisplayModeMulti).
				Sort(common.SortOrderDescending),
		)
	if len(metric.thresholds) > 0 {
		builder = builder.
			ColorScheme(
				dashboard.NewFieldColorBuilder().
					Mode(dashboard.FieldColorModeIdThresholds).
					SeriesBy(dashboard.FieldColorSeriesByModeLast),
			).
			ThresholdsStyle(
				common.NewGraphThresholdsStyleConfigBuilder().
					Mode(common.GraphThresholdsStyleModeDashed),
			).
			GradientMode(common.GraphGradientModeScheme)
	}
	return builder
}

func buildTemperatureTimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricTemperature, dashboard.GridPos{H: 13, W: 12, X: 0, Y: 10})
}

func buildHumidityTimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricHumidity, dashboard.GridPos{H: 13, W: 12, X: 12, Y: 10})
}

// ============================================================================
// Timeseries Panels - CO2
// ============================================================================

func buildCO2TimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricCO2, dashboard.GridPos{H: 13, W: 24, X: 0, Y: 24}).
		// Atmospheric CO2 levels are typically around 430ppm
		AxisSoftMin(430)
}

// ============================================================================
// Timeseries Panels - Particle Pollution
// ============================================================================

func buildPM03TimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricPM03, dashboard.GridPos{H: 14, W: 12, X: 0, Y: 38})
}

func buildPM1TimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricPM1, dashboard.GridPos{H: 14, W: 12, X: 12, Y: 38})
}

func buildPM25TimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricPM25, dashboard.GridPos{H: 14, W: 12, X: 0, Y: 52})
}

func buildPM10TimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricPM10, dashboard.GridPos{H: 14, W: 12, X: 12, Y: 52})
}

// ============================================================================
// Timeseries Panels - TVOC / NOx
// ============================================================================

func buildTVOCTimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricTVOC, dashboard.GridPos{H: 14, W: 12, X: 0, Y: 67})
}

func buildNOxTimeseriesPanel() *timeseries.PanelBuilder {
	return buildTimeseriesPanel(metricNOx, dashboard.GridPos{H: 14, W: 12, X: 12, Y: 67})
}

// ============================================================================
// Helpers
// ============================================================================

func buildThresholds(thresholds thresholdDefinition) *dashboard.ThresholdsConfigBuilder {
	thresholdsSteps := make([]dashboard.Threshold, 0, len(thresholds))
	for value, color := range thresholds {
		thresholdsSteps = append(thresholdsSteps, dashboard.Threshold{
			Color: color,
			Value: cog.ToPtr(value),
		})
	}
	sort.Slice(thresholdsSteps, func(i, j int) bool {
		return *thresholdsSteps[i].Value < *thresholdsSteps[j].Value
	})
	return dashboard.NewThresholdsConfigBuilder().
		Mode(dashboard.ThresholdsModeAbsolute).
		Steps(thresholdsSteps)
}

func prometheusDatasourceRef() dashboard.DataSourceRef {
	return dashboard.DataSourceRef{
		Type: cog.ToPtr("prometheus"),
	}
}

func prometheusQuery(expr, legendFormat string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(expr).
		LegendFormat(legendFormat).
		EditorMode(prometheus.QueryEditorModeCode)
}
