package flags

import (
	"flag"
	"os"
	"strconv"
)

var ServerEndPoint string

func ParseServerFlags() {
	flag.StringVar(&ServerEndPoint, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envServerEndPoint := os.Getenv("ADDRESS"); envServerEndPoint != "" {
		ServerEndPoint = envServerEndPoint
	}
}

var (
	AgentEndPoint  string
	ReportInterval int64
	PollInterval   int64
)

func ParseAgentFlag() error {
	flag.StringVar(&AgentEndPoint, "a", "localhost:8080", "address and port to run agent")
	flag.Int64Var(&ReportInterval, "r", 10, "frequency of report metrics")
	flag.Int64Var(&PollInterval, "p", 2, "the frequency of the metric survey")
	flag.Parse()

	var err error

	if envAgentEndPoint := os.Getenv("ADDRESS"); envAgentEndPoint != "" {
		AgentEndPoint = envAgentEndPoint
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		ReportInterval, err = strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			return err
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			return err
		}
	}
	return nil
}
