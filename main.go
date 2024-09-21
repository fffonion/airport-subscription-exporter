package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/fffonion/airport-subscription-exporter/exporter"
)

func main() {
	var metricsAddr = flag.String("metrics.listen-addr", ":9991", "listen address for airport subscription exporter")

	flag.Parse()
	s := exporter.NewHttpServer()
	log.Infof("Accepting Prometheus Requests on %s", *metricsAddr)
	log.Fatal(http.ListenAndServe(*metricsAddr, s))
}
