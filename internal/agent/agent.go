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

	body, err := m.Compress(metrics)
	if err != nil {
		log.Printf("error compress procedure. Err : %s", err.Error())
		return
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post("http://" + url + "/updates/")

	if err != nil {
		log.Printf(err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("unexpected status code. Expected code 200, got %d.", resp.StatusCode())
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
