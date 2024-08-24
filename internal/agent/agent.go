package agent

import (
	"encoding/json"
	"errors"
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func postMetrics(url string, metrics []m.Metric) error {
	client := resty.New()
	url, _ = strings.CutPrefix(url, "http://")

	for _, metric := range metrics {
		metricJSON, err := json.Marshal(metric)
		if err != nil {
			log.Printf("error serializing metric: %s. Name metric: %s", err, metric.ID)
			return err
		}
		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetBody(metricJSON).Post("http://" + url + "/update/")

		if err != nil {
			log.Printf("error send request: %s. Name metric: %s", err, metric.ID)
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("unexpected status code. Expected code 200, got %d. Name metric: %s", resp.StatusCode(), metric.ID)
			return errors.New("unexpected status code")
		}
	}
	m.PollCount = 0
	return nil
}

func Start() error {
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
			metrics = append(metrics, m.Metric{MType: "counter", ID: "PollCount", Delta: &m.PollCount})
			err := postMetrics(f.AgentEndPoint, metrics)
			if err != nil {
				return err
			}
			metrics = make([]m.Metric, 0)
		}
	}
}
