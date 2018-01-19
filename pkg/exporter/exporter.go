package exporter

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ifosch/cloudinary-exporter/pkg/cloudinary"
)

const ns = "cloudinary"

type ReportDesc struct {
	Name  string
	Desc  string
	Value float64
}

var ReportDescs = []ReportDesc{
	ReportDesc{Name: "transformations_usage", Desc: "Transformation usage", Value: 1.0},
	ReportDesc{Name: "transformations_limit", Desc: "Transformation limit", Value: 1.0},
	ReportDesc{Name: "transformations_used_percent", Desc: "Transformation used percent", Value: 1.0},
	ReportDesc{Name: "objects_usage", Desc: "Object usage", Value: 1.0},
	ReportDesc{Name: "objects_limit", Desc: "Object limit", Value: 1.0},
	ReportDesc{Name: "objects_used_percent", Desc: "Object used percent", Value: 1.0},
	ReportDesc{Name: "bandwidth_usage", Desc: "Bandwidth usage", Value: 1.0},
	ReportDesc{Name: "bandwidth_limit", Desc: "Bandwidth limit", Value: 1.0},
	ReportDesc{Name: "bandwidth_used_percent", Desc: "Bandwidth used percent", Value: 1.0},
	ReportDesc{Name: "storage_usage", Desc: "Storage usage", Value: 1.0},
	ReportDesc{Name: "storage_limit", Desc: "Storage limit", Value: 1.0},
	ReportDesc{Name: "storage_used_percent", Desc: "Storage used percent", Value: 1.0},
	ReportDesc{Name: "requests", Desc: "Requests", Value: 1.0},
	ReportDesc{Name: "resources", Desc: "Resources", Value: 1.0},
	ReportDesc{Name: "derived_resources", Desc: "Derived resources", Value: 1.0},
}

type Exporter struct {
	metrics []prometheus.Gauge
}

func NewExporter() (*Exporter, error) {
	metrics := []prometheus.Gauge{}
	for _, desc := range ReportDescs {
		metricDesc := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        desc.Name,
				Help:        desc.Desc,
				ConstLabels: nil,
			},
		)
		metrics = append(metrics, metricDesc)
	}
	return &Exporter{
		metrics: metrics,
	}, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
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

	e.metrics[0].Set(float64(usageReport.Transformations.Usage))
	e.metrics[1].Set(float64(usageReport.Transformations.Limit))
	e.metrics[2].Set(float64(usageReport.Transformations.UsedPercent))
	e.metrics[3].Set(float64(usageReport.Objects.Usage))
	e.metrics[4].Set(float64(usageReport.Objects.Limit))
	e.metrics[5].Set(float64(usageReport.Objects.UsedPercent))
	e.metrics[6].Set(float64(usageReport.Bandwidth.Usage))
	e.metrics[7].Set(float64(usageReport.Bandwidth.Limit))
	e.metrics[8].Set(float64(usageReport.Bandwidth.UsedPercent))
	e.metrics[9].Set(float64(usageReport.Storage.Usage))
	e.metrics[10].Set(float64(usageReport.Storage.Limit))
	e.metrics[11].Set(float64(usageReport.Storage.UsedPercent))
	e.metrics[12].Set(float64(usageReport.Requests))
	e.metrics[13].Set(float64(usageReport.Resources))
	e.metrics[14].Set(float64(usageReport.DerivedResources))
	log.Println(*usageReport)

	return nil
}
