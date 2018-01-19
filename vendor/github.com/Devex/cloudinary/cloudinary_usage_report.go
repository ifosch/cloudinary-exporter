package cloudinary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Devex/heracles"
)

const ns = "cloudinary"

var cloudinaryMetrics = heracles.GaugeMetrics{
	1:  heracles.NewGaugeMetrics(ns, "transformations_usage", "Transformations usage", nil),
	2:  heracles.NewGaugeMetrics(ns, "transformations_limit", "Transformations limit", nil),
	3:  heracles.NewGaugeMetrics(ns, "transformations_used_percent", "Transformations used percent", nil),
	4:  heracles.NewGaugeMetrics(ns, "objects_usage", "Objects usage", nil),
	5:  heracles.NewGaugeMetrics(ns, "objects_limit", "Objects limit", nil),
	6:  heracles.NewGaugeMetrics(ns, "objects_used_percent", "Objects used percent", nil),
	7:  heracles.NewGaugeMetrics(ns, "bandwidth_usage", "Bandwidth usage", nil),
	8:  heracles.NewGaugeMetrics(ns, "bandwidth_limit", "Bandwidth limit", nil),
	9:  heracles.NewGaugeMetrics(ns, "bandwidth_used_percent", "Bandwidth used percent", nil),
	10: heracles.NewGaugeMetrics(ns, "storage_usage", "Storage usage", nil),
	11: heracles.NewGaugeMetrics(ns, "storage_limit", "Storage limit", nil),
	12: heracles.NewGaugeMetrics(ns, "storage_used_percent", "Storage used percent", nil),
	13: heracles.NewGaugeMetrics(ns, "requests", "Requests", nil),
	14: heracles.NewGaugeMetrics(ns, "resources", "Resources", nil),
	15: heracles.NewGaugeMetrics(ns, "derived_resources", "Derived resources", nil),
}

type UsageInfo struct {
	Usage       int64   `json:"usage"`
	Limit       int64   `json:"limit"`
	UsedPercent float64 `json:"used_percent"`
}

type UsageReport struct {
	Plan             string    `json:"plan"`
	LastUpdate       string    `json:"last_updated"`
	Transformations  UsageInfo `json:"transformations"`
	Objects          UsageInfo `json:"objects"`
	Bandwidth        UsageInfo `json:"bandwidth"`
	Storage          UsageInfo `json:"storage"`
	Requests         int64     `json:"requests"`
	Resources        int64     `json:"resources"`
	DerivedResources int64     `json:"derived_resources"`
}

type Exporter struct {
	metrics heracles.GaugeMetrics
}

func NewExporter() (*Exporter, error) {
	return &Exporter{
		metrics: cloudinaryMetrics,
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collect called...")
	err := e.fetch()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range e.metrics {
		m.Collect(ch)
	}
}

func (e *Exporter) fetch() (err error) {
	key, secret, cloud_name, err := getCredentials()
	if err != nil {
		return err
	}

	rs, err := http.Get(
		fmt.Sprintf(
			"https://%s:%s@api.cloudinary.com/v1_1/%s/usage",
			key,
			secret,
			cloud_name,
		),
	)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err
	}

	usageReport := new(UsageReport)
	err = json.Unmarshal(bodyBytes, &usageReport)
	if err != nil {
		return err
	}
	e.metrics[1].With(nil).Set(float64(usageReport.Transformations.Usage))
	e.metrics[2].With(nil).Set(float64(usageReport.Transformations.Limit))
	e.metrics[3].With(nil).Set(float64(usageReport.Transformations.UsedPercent))
	e.metrics[4].With(nil).Set(float64(usageReport.Objects.Usage))
	e.metrics[5].With(nil).Set(float64(usageReport.Objects.Limit))
	e.metrics[6].With(nil).Set(float64(usageReport.Objects.UsedPercent))
	e.metrics[7].With(nil).Set(float64(usageReport.Bandwidth.Usage))
	e.metrics[8].With(nil).Set(float64(usageReport.Bandwidth.Limit))
	e.metrics[9].With(nil).Set(float64(usageReport.Bandwidth.UsedPercent))
	e.metrics[10].With(nil).Set(float64(usageReport.Storage.Usage))
	e.metrics[11].With(nil).Set(float64(usageReport.Storage.Limit))
	e.metrics[12].With(nil).Set(float64(usageReport.Storage.UsedPercent))
	e.metrics[13].With(nil).Set(float64(usageReport.Requests))
	e.metrics[14].With(nil).Set(float64(usageReport.Resources))
	e.metrics[15].With(nil).Set(float64(usageReport.DerivedResources))
	log.Println(usageReport)

	return nil
}

func getCredentials() (key, secret, cloud_name string, err error) {
	key = os.Getenv("CLOUDINARY_KEY")
	secret = os.Getenv("CLOUDINARY_SECRET")
	cloud_name = os.Getenv("CLOUDINARY_CLOUD_NAME")

	if key == "" || secret == "" || cloud_name == "" {
		err = errors.New("No credentials defined")
	}
	return key, secret, cloud_name, err
}
