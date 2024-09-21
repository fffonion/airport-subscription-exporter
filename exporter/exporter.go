package exporter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	subscription_url string

	metricsUp,
	metricsRefreshed,
	metricsUpload,
	metricsDownload,
	metricsTotal,
	metricsExpire *prometheus.Desc
}

type ExporterTarget struct {
	URL string
}

func NewExporter(t *ExporterTarget) *Exporter {
	var (
		constLabels = prometheus.Labels{}
		labelNames  = []string{}
	)

	e := &Exporter{
		subscription_url: t.URL,
		metricsUp: prometheus.NewDesc("airport_online",
			"Airpot collector online",
			nil, constLabels,
		),

		metricsRefreshed: prometheus.NewDesc("airport_last_refreshed",
			"Airport subscription info last refreshed time",
			labelNames, constLabels,
		),

		metricsUpload: prometheus.NewDesc("airport_userinfo_upload",
			"Airport upload traffic",
			labelNames, constLabels,
		),

		metricsDownload: prometheus.NewDesc("airport_userinfo_download",
			"Airport download traffic",
			labelNames, constLabels,
		),

		metricsTotal: prometheus.NewDesc("airport_userinfo_total",
			"Airport total subscription traffic",
			labelNames, constLabels,
		),

		metricsExpire: prometheus.NewDesc("airport_userinfo_expire",
			"Airport subscription expire time",
			labelNames, constLabels,
		),
	}
	return e
}

func (k *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- k.metricsUp
	ch <- k.metricsRefreshed
	ch <- k.metricsUpload
	ch <- k.metricsDownload
	ch <- k.metricsTotal
	ch <- k.metricsExpire

}

func (k *Exporter) Collect(ch chan<- prometheus.Metric) {
	sub, err := parse(k.subscription_url)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(k.metricsUp, prometheus.GaugeValue,
			0)
		log.Errorln("error collecting", k.subscription_url, ":", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(k.metricsRefreshed, prometheus.GaugeValue,
		float64(time.Now().Unix()))

	ch <- prometheus.MustNewConstMetric(k.metricsUpload, prometheus.CounterValue,
		float64(sub.upload))

	ch <- prometheus.MustNewConstMetric(k.metricsDownload, prometheus.CounterValue,
		float64(sub.download))

	ch <- prometheus.MustNewConstMetric(k.metricsTotal, prometheus.CounterValue,
		float64(sub.total))

	ch <- prometheus.MustNewConstMetric(k.metricsExpire, prometheus.GaugeValue,
		float64(sub.expire))
}
