package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/avast/retry-go"
	m "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent/metrics"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func postMetrics(url string, metrics []m.Metric) error {
	client := resty.New()
	url, _ = strings.CutPrefix(url, "http://")

	body, err := m.Compress(metrics)
	if err != nil {
		log.Printf("error compress procedure. Err : %s", err.Error())
		return err
	}
	uncompressedBody, err := m.GetMetricsAsBody(metrics)
	if err != nil {
		log.Printf("error convert procedure. Err : %s", err.Error())
		return err
	}
	hash := getHash(uncompressedBody, f.AgentKey)

	err = retry.Do(func() error {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("HashSHA256", hash).
			SetBody(body).
			Post("http://" + url + "/updates/")

		if err != nil {
			log.Println(err.Error())
			if !utils.IsRetriableError(err) {
				return retry.Unrecoverable(err)
			}
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("unexpected status code. Expected code 200, got %d.", resp.StatusCode())
		}

		m.PollCount = 0

		return nil
	},
		retry.Attempts(4),
		retry.DelayType(utils.RetryDelay))
	return err
}

func getHash(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
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
			err := postMetrics(f.AgentEndPoint, metrics)
			if err != nil {
				log.Printf("%s - attempts to send metrics failed", err.Error())
			}
			metrics = make([]m.Metric, 0)
		}
	}
}
