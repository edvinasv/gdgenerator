package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type DashboardDataTimeRange struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type DataRowColumnGroupItem struct {
	Title         string   `yaml:"title"`
	Descr         []string `yaml:"description"`
	Expr          string   `yaml:"expr"`
	Link          string   `yaml:"link"`
	Min           float64  `yaml:"min"`
	Max           float64  `yaml:"max"`
	DatasourceUid string   `yaml:"datasourceUid"`
	Type          string   `yaml:"type"`
	Service       string   `yaml:"service"`
	Disabled      bool     `yaml:"disabled"`
}

type DataRowColumnGroup struct {
	Title         string                   `yaml:"title"`
	Width         uint32                   `yaml:"width"`
	Height        uint32                   `yaml:"height"`
	Min           float64                  `yaml:"min"`
	Max           float64                  `yaml:"max"`
	Spacing       uint32                   `yaml:"spacing"`
	Columns       uint32                   `yaml:"columns"`
	DatasourceUid string                   `yaml:"datasourceUid"`
	Type          string                   `yaml:"type"`
	Service       string                   `yaml:"service"`
	Disabled      bool                     `yaml:"disabled"`
	Items         []DataRowColumnGroupItem `yaml:"items"`
}

type DataRowColumn struct {
	Title         string               `yaml:"title"`
	DatasourceUid string               `yaml:"datasourceUid"`
	Type          string               `yaml:"type"`
	Service       string               `yaml:"service"`
	Disabled      bool                 `yaml:"disabled"`
	Groups        []DataRowColumnGroup `yaml:"groups"`
}

type DataRow struct {
	Title         string          `yaml:"title"`
	Columns       []DataRowColumn `yaml:"columns"`
	DatasourceUid string          `yaml:"datasourceUid"`
	Type          string          `yaml:"type"`
	Service       string          `yaml:"service"`
	Disabled      bool            `yaml:"disabled"`
}

type DashboardData struct {
	Name          string                 `yaml:"title"`
	Uid           string                 `yaml:"uid"`
	FolderUid     string                 `yaml:"folderUid"`
	TimeRange     DashboardDataTimeRange `yaml:"timeRange"`
	Refresh       string                 `yaml:"refresh"`
	DatasourceUid string                 `yaml:"datasourceUid"`
	Type          string                 `yaml:"type"`
	Service       string                 `yaml:"service"`
	Disabled      bool                   `yaml:"disabled"`
	Banner        string                 `yaml:"banner"`
	DryRun        bool                   `yaml:"dryRun"`
	Rows          []DataRow              `yaml:"rows"`
}

type AlertLabels struct {
	Severity string `yaml:"severity,omitempty"`
	Service  string `yaml:"service,omitempty"`
}

type AlertAnnotations struct {
	Summary string `yaml:"summary"`
}

type AlertRule struct {
	Alert                 string           `yaml:"alert"`
	Expr                  string           `yaml:"expr"`
	ForDuration           string           `yaml:"for,omitempty"`
	KeepFiringForDuration string           `yaml:"keep_firing_for,omitempty"`
	Labels                AlertLabels      `yaml:"labels"`
	Annotations           AlertAnnotations `yaml:"annotations"`
}

type AlertGroup struct {
	Name     string      `yaml:"name"`
	Interval string      `yaml:"interval,omitempty"`
	Rules    []AlertRule `yaml:"rules"`
}

type AlertGroups struct {
	Groups []AlertGroup `yaml:"groups"`
}

func main() {

	manifests := false
	manifestsDirectory := "./resources"
	alerts := false
	alertsDirectory := "./alerts"
	configFile := "examples/generic.yaml"

	var dashboardData DashboardData

	flag.BoolVar(&manifests, "manifests", manifests, "Generate a dashboard manifest and write it to disk")
	flag.StringVar(&manifestsDirectory, "manifests-directory", manifestsDirectory, "Directory in which the manifests will be generated")
	flag.BoolVar(&alerts, "alerts", alerts, "Generate alert configuration and write it to disk")
	flag.StringVar(&alertsDirectory, "alerts-directory", alertsDirectory, "Alerts configuration directory")
	flag.StringVar(&configFile, "config", configFile, "Dashboard configuration file")
	flag.Parse()

	dashboardData.getConf(configFile)

	dashboard, err := generateDashboard(dashboardData).Build()
	if err != nil {
		log.Fatal(err)
	}

	if manifests {
		if err := generateDashboardManifest(dashboardData.FolderUid, manifestsDirectory, dashboard); err != nil {
			log.Fatal(err)
		}
		return
	}

	if alerts {
		alertGroups := generateAlerts(dashboardData)
		if err := generateAlertGroupFile(alertsDirectory, dashboard, alertGroups); err != nil {
			log.Fatal(err)
		}
		return
	}

	printDashboard(dashboard)

}

func (c *DashboardData) getConf(configFile string) *DashboardData {

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func generateDashboardManifest(folderUid string, outputDir string, dashboard dashboard.Dashboard) error {

	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return err
	}

	manifest := DashboardManifest(folderUid, dashboard)
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	filename := *dashboard.Uid + ".json"
	if err := os.WriteFile(filepath.Join(outputDir, filename), manifestJSON, 0666); err != nil {
		return err
	}

	return nil
}

func generateAlertGroupFile(outputDir string, dashboard dashboard.Dashboard, alertGroups *AlertGroups) error {
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return err
	}

	alertConfig, err := yaml.Marshal(alertGroups)
	if err != nil {
		return err
	}

	filename := *dashboard.Uid + ".rules.yaml"
	if err := os.WriteFile(filepath.Join(outputDir, filename), alertConfig, 0666); err != nil {
		return err
	}

	return nil

}

func printDashboard(dashboard dashboard.Dashboard) {
	manifest := DashboardManifest("", dashboard)
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(manifestJSON))
}
