package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ifosch/cloudinary-exporter/pkg/exporter"
)

func getListenAddress() string {
	port := flag.Int("p", 9101, "Port to use for listen")
	address := flag.String("a", "", "Address to use for listen")
	flag.Parse()

	return fmt.Sprintf("%v:%d", *address, *port)
}

func main() {
	listenAddress := getListenAddress()
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
