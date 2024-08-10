package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	agentEndPoint  string
	reportInterval int64
	pollInterval   int64
)

func parseFlag() {
	flag.StringVar(&agentEndPoint, "a", "localhost:8080", "address and port to run agent")
	flag.Int64Var(&reportInterval, "r", 10, "frequency of report metrics")
	flag.Int64Var(&pollInterval, "p", 2, "the frequency of the metric survey")
	flag.Parse()

	if envAgentEndPoint := os.Getenv("ADDRESS"); envAgentEndPoint != "" {
		agentEndPoint = envAgentEndPoint
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, _ = strconv.ParseInt(envReportInterval, 10, 64)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, _ = strconv.ParseInt(envPollInterval, 10, 64)
	}
}
