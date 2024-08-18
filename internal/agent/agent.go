package agent

import (
	"errors"
	"fmt"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/metrics"
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
		resp, err := client.R().SetPathParams(map[string]string{
			"type":  metric.Type,
			"name":  metric.Name,
			"value": fmt.Sprintf("%v", metric.Value),
		}).SetHeader("Content-Type", "text/plain").Post("http://" + url + "/update/{type}/{name}/{value}")

		if err != nil {
			log.Printf("error send request: %s. Name metric: %s", err, metric.Name)
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("unexpected status code. Expected code 200, got %d. Name metric: %s", resp.StatusCode(), metric.Name)
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
			metrics = append(metrics, m.GetAllMetrics()...)

		case <-reportTicker.C:
			err := postMetrics(f.AgentEndPoint, metrics)
			if err != nil {
				return err
			}
			metrics = make([]m.Metric, 0)
		}
	}
}
