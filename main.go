package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ifosch/cloudinary-exporter/pkg/exporter"
)

func main() {
	listenAddress := ":9101"
	log.Println("Starting cloudinary-exporter")

	exporter, err := exporter.NewExporter()
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(exporter)

	log.Println("Listening on", listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
                <head><title>Cloudinary Exporter</title></head>
                <body>
                <h1>Cloudinary Exporter</h1>
                <p><a href='/metrics'>Metrics</a></p>
                </body>
                </html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
