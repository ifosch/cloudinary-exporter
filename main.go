package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ifosch/cloudinary-exporter/pkg/cloudinary"
	"github.com/ifosch/cloudinary-exporter/pkg/exporter"
)

func getListenAddress() string {
	port := flag.Int("p", 9101, "Port to use for listen")
	address := flag.String("a", "", "Address to use for listen")
	flag.Parse()

	return fmt.Sprintf("%v:%d", *address, *port)
}

var l *log.Logger

func main() {
	listenAddress := getListenAddress()
	l := log.New(os.Stderr, "", 1)
	l.Println("Starting cloudinary-exporter")

	err := cloudinary.NewCredentials(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_KEY"),
		os.Getenv("CLOUDINARY_SECRET"),
	)
	if err != nil {
		l.Fatal(err)
	}

	exporter, err := exporter.NewExporter(l)
	if err != nil {
		l.Fatal(err)
	}

	prometheus.MustRegister(exporter)

	l.Println("Listening on", listenAddress)
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
	l.Fatal(http.ListenAndServe(listenAddress, nil))
}
