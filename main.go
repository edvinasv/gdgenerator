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
	Disabled      bool                     `yaml:"disabled"`
	Items         []DataRowColumnGroupItem `yaml:"items"`
}

type DataRowColumn struct {
	Title         string               `yaml:"title"`
	DatasourceUid string               `yaml:"datasourceUid"`
	Type          string               `yaml:"type"`
	Disabled      bool                 `yaml:"disabled"`
	Groups        []DataRowColumnGroup `yaml:"groups"`
}

type DataRow struct {
	Title         string          `yaml:"title"`
	Columns       []DataRowColumn `yaml:"columns"`
	DatasourceUid string          `yaml:"datasourceUid"`
	Type          string          `yaml:"type"`
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
	Disabled      bool                   `yaml:"disabled"`
	Banner        string                 `yaml:"banner"`
	DryRun        bool                   `yaml:"dryRun"`
	Rows          []DataRow              `yaml:"rows"`
}

func main() {

	manifests := false
	manifestsDirectory := "./resources"
	configFile := "examples/generic.yaml"

	var dashboardData DashboardData

	flag.BoolVar(&manifests, "manifests", manifests, "Generate a dashboard manifest for the test dashboard and write it to disk")
	flag.StringVar(&manifestsDirectory, "manifests-directory", manifestsDirectory, "Directory in which the manifests will be generated")
	flag.StringVar(&configFile, "config", configFile, "Dashboard configuration file")
	flag.Parse()

	dashboardData.getConf(configFile)

	genericDashboard, err := generateDashboard(dashboardData).Build()
	if err != nil {
		log.Fatal(err)
	}

	// Generate a dashboard manifest for the generic dashboard and write it to disk.
	if manifests {
		if err := generateManifest(dashboardData.FolderUid, manifestsDirectory, genericDashboard); err != nil {
			log.Fatal(err)
		}
		return
	}

	// By default: print the test dashboard to stdout.
	printDashboard(genericDashboard)
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

func generateManifest(folderUid string, outputDir string, dashboard dashboard.Dashboard) error {

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

func printDashboard(dashboard dashboard.Dashboard) {
	manifest := DashboardManifest("", dashboard)
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(manifestJSON))
}
