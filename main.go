package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/fffonion/airport-subscription-exporter/exporter"
)

func parseDuration(duration string) (int, error) {
	if len(duration) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	value, err := strconv.Atoi(duration[:len(duration)-1])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %v", err)
	}

	unit := duration[len(duration)-1:]
	switch unit {
	case "d":
		return value * 24 * 60 * 60, nil
	case "h":
		return value * 60 * 60, nil
	case "m":
		return value * 60, nil
	case "s":
		return value, nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %s", unit)
	}
}

func main() {
	var metricsAddr = flag.String("metrics.listen-addr", ":9991", "listen address for airport subscription exporter")
	var subscriptionUpdateInternal = flag.String("sub.update-interval", "1h", "how long should exporter actually update subscription info")

	parseDuration, err := parseDuration(*subscriptionUpdateInternal)
	if err != nil {
		log.Fatalf("failed to parse duration: %v", err)
		return
	}
	flag.Parse()
	s := exporter.NewHttpServer(parseDuration)
	log.Infof("Accepting Prometheus Requests on %s", *metricsAddr)
	log.Fatal(http.ListenAndServe(*metricsAddr, s))
}
