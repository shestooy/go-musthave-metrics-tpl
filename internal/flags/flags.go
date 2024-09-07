package flags

import (
	"flag"
	"os"
	"strconv"
)

var (
	ServerEndPoint  string
	LogLevel        string
	StorageInterval int64
	FileStoragePath string
	Restore         bool
	AddrDB          string
)

func ParseServerFlags() error {
	var err error

	flag.StringVar(&ServerEndPoint, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&LogLevel, "l", "info", "log level")
	flag.Int64Var(&StorageInterval, "i", 300, "the time interval in seconds for saving metrics to disk")
	flag.StringVar(&FileStoragePath, "f", "metric.txt", "the path to the file for storing metrics")
	flag.BoolVar(&Restore, "r", true, "whether to load saved metrics at startup")
	flag.StringVar(&AddrDB, "d", "", "the address of the database")
	flag.Parse()

	if envServerEndPoint := os.Getenv("ADDRESS"); envServerEndPoint != "" {
		ServerEndPoint = envServerEndPoint
	}
	if envFlagLogLevel := os.Getenv("LOG_LEVEL"); envFlagLogLevel != "" {
		LogLevel = envFlagLogLevel
	}
	if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
		StorageInterval, err = strconv.ParseInt(envStorageInterval, 10, 64)
		if err != nil {
			return err
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return err
		}
	}
	if envDBDsn := os.Getenv("DATABASE_DSN"); envDBDsn != "" {
		AddrDB = envDBDsn
	}
	return nil
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
