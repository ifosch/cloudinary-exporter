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
	Value func(cloudinary.UsageReport) float64
}

var ReportDescs = []ReportDesc{
	ReportDesc{Name: "transformations_usage", Desc: "Transformation usage", Value: cloudinary.TransformationUsage},
	ReportDesc{Name: "transformations_limit", Desc: "Transformation limit", Value: cloudinary.TransformationLimit},
	ReportDesc{Name: "transformations_used_percent", Desc: "Transformation used percent", Value: cloudinary.TransformationUsedPercent},
	ReportDesc{Name: "objects_usage", Desc: "Object usage", Value: cloudinary.ObjectsUsage},
	ReportDesc{Name: "objects_limit", Desc: "Object limit", Value: cloudinary.ObjectsLimit},
	ReportDesc{Name: "objects_used_percent", Desc: "Object used percent", Value: cloudinary.ObjectsUsedPercent},
	ReportDesc{Name: "bandwidth_usage", Desc: "Bandwidth usage", Value: cloudinary.BandwidthUsage},
	ReportDesc{Name: "bandwidth_limit", Desc: "Bandwidth limit", Value: cloudinary.BandwidthLimit},
	ReportDesc{Name: "bandwidth_used_percent", Desc: "Bandwidth used percent", Value: cloudinary.BandwidthUsedPercent},
	ReportDesc{Name: "storage_usage", Desc: "Storage usage", Value: cloudinary.StorageUsage},
	ReportDesc{Name: "storage_limit", Desc: "Storage limit", Value: cloudinary.StorageLimit},
	ReportDesc{Name: "storage_used_percent", Desc: "Storage used percent", Value: cloudinary.StorageUsedPercent},
	ReportDesc{Name: "requests", Desc: "Requests", Value: cloudinary.Requests},
	ReportDesc{Name: "resources", Desc: "Resources", Value: cloudinary.Resources},
	ReportDesc{Name: "derived_resources", Desc: "Derived resources", Value: cloudinary.DerivedResources},
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

	report, err := cloudinary.GetUsageReport(req)
	if err != nil {
		return err
	}

	for i, metric := range ReportDescs {
		if metric.Value != nil {
			e.metrics[i].Set(metric.Value(*report))
		}
	}

	return nil
}
