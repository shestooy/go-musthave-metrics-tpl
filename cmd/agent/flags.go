package main

import "flag"

var (
	agentEndPoint  string
	reportInterval int64
	pollInterval   int64
)

func parseFlag() {
	flag.StringVar(&agentEndPoint, "a", "http://localhost:8080", "address and port to run server")
	flag.Int64Var(&reportInterval, "r", 10, "frequency of report metrics")
	flag.Int64Var(&pollInterval, "p", 2, "the frequency of the metric survey")

	flag.Parse()
}
