package agent

import (
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func postMetrics(url string, metrics []m.Metric) {
	client := resty.New()
	url, _ = strings.CutPrefix(url, "http://")

	for _, metric := range metrics {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(metric.Compress()).
			Post("http://" + url + "/update/")

		if err != nil {
			log.Printf("error send request: %s. Name metric: %s", err, metric.ID)
			continue
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("unexpected status code. Expected code 200, got %d. Name metric: %s", resp.StatusCode(), metric.ID)
			continue
		}
	}
	m.PollCount = 0
}

func Start() {
	pollTicker := time.NewTicker(time.Duration(f.PollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(f.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	metrics := make([]m.Metric, 0)

	for {
		select {
		case <-pollTicker.C:
			m.PollCount++
			metrics = append(metrics, m.GetAllMetrics()...)

		case <-reportTicker.C:
			metrics = append(metrics, m.Metric{MType: m.Counter, ID: "PollCount", Delta: &m.PollCount})
			postMetrics(f.AgentEndPoint, metrics)
			metrics = make([]m.Metric, 0)
		}
	}
}
