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
	ReportDesc{Name: "transformations_usage_amount", Desc: "Number of used transformations in the last 30 days", Value: cloudinary.TransformationUsage},
	ReportDesc{Name: "transformations_limit_amount", Desc: "Limit of transformations allowed in the last 30 days", Value: cloudinary.TransformationLimit},
	ReportDesc{Name: "transformations_usage_ratio", Desc: "Ratio of used transformations over corresponding limit", Value: cloudinary.TransformationUsageRatio},
	ReportDesc{Name: "objects_usage_amount", Desc: "Number of used objects in the last 30 days", Value: cloudinary.ObjectsUsage},
	ReportDesc{Name: "objects_limit_amount", Desc: "Limit of objects allowed in the last 30 days", Value: cloudinary.ObjectsLimit},
	ReportDesc{Name: "objects_usage_ratio", Desc: "Ratio of used objects over corresponding limit", Value: cloudinary.ObjectsUsageRatio},
	ReportDesc{Name: "bandwidth_usage_bytes", Desc: "Bytes used in bandwidth in the last 30 days", Value: cloudinary.BandwidthUsage},
	ReportDesc{Name: "bandwidth_limit_bytes", Desc: "Limit of bytes in bandwidth to use in the last 30 days", Value: cloudinary.BandwidthLimit},
	ReportDesc{Name: "bandwidth_usage_ratio", Desc: "Ratio of used bytes in bandwidth over corresponding limit", Value: cloudinary.BandwidthUsageRatio},
	ReportDesc{Name: "storage_usage_bytes", Desc: "Bytes of storage used in the last 30 days", Value: cloudinary.StorageUsage},
	ReportDesc{Name: "storage_limit_bytes", Desc: "Limit of storage used allowed in the last 30 days", Value: cloudinary.StorageLimit},
	ReportDesc{Name: "storage_usage_ratio", Desc: "Ratio of storage bytes used over corresponding limit", Value: cloudinary.StorageUsageRatio},
	ReportDesc{Name: "requests_amount", Desc: "Number of requests done to Cloudinary", Value: cloudinary.Requests},
	ReportDesc{Name: "resources_amount", Desc: "Number of resources in Cloudinary", Value: cloudinary.Resources},
	ReportDesc{Name: "derived_resources_amount", Desc: "Number of derived resources in Cloudinary", Value: cloudinary.DerivedResources},
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
