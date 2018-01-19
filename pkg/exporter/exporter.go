package exporter

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ifosch/cloudinary-exporter/pkg/cloudinary"
)

type GaugeMetrics map[int]*prometheus.GaugeVec

func NewGaugeMetrics(
	namespace, name, docString string,
	labels prometheus.Labels,
) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        name,
			Help:        docString,
			ConstLabels: labels,
		},
		[]string{},
	)
}

const ns = "cloudinary"

var cloudinaryMetrics = GaugeMetrics{
	1:  NewGaugeMetrics(ns, "transformations_usage", "Transformations usage", nil),
	2:  NewGaugeMetrics(ns, "transformations_limit", "Transformations limit", nil),
	3:  NewGaugeMetrics(ns, "transformations_used_percent", "Transformations used percent", nil),
	4:  NewGaugeMetrics(ns, "objects_usage", "Objects usage", nil),
	5:  NewGaugeMetrics(ns, "objects_limit", "Objects limit", nil),
	6:  NewGaugeMetrics(ns, "objects_used_percent", "Objects used percent", nil),
	7:  NewGaugeMetrics(ns, "bandwidth_usage", "Bandwidth usage", nil),
	8:  NewGaugeMetrics(ns, "bandwidth_limit", "Bandwidth limit", nil),
	9:  NewGaugeMetrics(ns, "bandwidth_used_percent", "Bandwidth used percent", nil),
	10: NewGaugeMetrics(ns, "storage_usage", "Storage usage", nil),
	11: NewGaugeMetrics(ns, "storage_limit", "Storage limit", nil),
	12: NewGaugeMetrics(ns, "storage_used_percent", "Storage used percent", nil),
	13: NewGaugeMetrics(ns, "requests", "Requests", nil),
	14: NewGaugeMetrics(ns, "resources", "Resources", nil),
	15: NewGaugeMetrics(ns, "derived_resources", "Derived resources", nil),
}

type Exporter struct {
	metrics GaugeMetrics
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
	req, err := cloudinary.GetRequest()
	if err != nil {
		return err
	}

	usageReport, err := cloudinary.GetUsageReport(req)
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
	log.Println(*usageReport)

	return nil
}
