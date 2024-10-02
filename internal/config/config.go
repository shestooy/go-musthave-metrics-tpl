package config

import (
	"flag"
	"os"
	"strconv"
)

type ServerCfg struct {
	ServerEndPoint  string
	LogLevel        string
	FileStoragePath string
	Restore         bool
	AddrDB          string
	ServerKey       string
	StorageInterval int64
}

func GetServerCfg() (*ServerCfg, error) {
	var err error
	cfg := &ServerCfg{}
	flag.StringVar(&cfg.ServerEndPoint, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.Int64Var(&cfg.StorageInterval, "i", 300, "the time interval in seconds for saving metrics to disk")
	flag.StringVar(&cfg.FileStoragePath, "f", "metric.txt", "the path to the file for storing metrics")
	flag.BoolVar(&cfg.Restore, "r", true, "whether to load saved metrics at startup")
	flag.StringVar(&cfg.AddrDB, "d", "", "the address of the database")
	flag.StringVar(&cfg.ServerKey, "k", "", "the server key for HashSHA256")
	flag.Parse()

	if envServerEndPoint := os.Getenv("ADDRESS"); envServerEndPoint != "" {
		cfg.ServerEndPoint = envServerEndPoint
	}
	if envFlagLogLevel := os.Getenv("LOG_LEVEL"); envFlagLogLevel != "" {
		cfg.LogLevel = envFlagLogLevel
	}
	if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
		cfg.StorageInterval, err = strconv.ParseInt(envStorageInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		cfg.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
	}
	if envDBDsn := os.Getenv("DATABASE_DSN"); envDBDsn != "" {
		cfg.AddrDB = envDBDsn
	}
	if envServerKey := os.Getenv("KEY"); envServerKey != "" {
		cfg.ServerKey = envServerKey
	}
	return cfg, nil
}

type AgentCfg struct {
	ReportInterval int64
	PollInterval   int64
	RateLimit      int64
	AgentKey       string
	AgentEndPoint  string
}

func GetAgentCfg() (*AgentCfg, error) {
	cfg := &AgentCfg{}
	flag.StringVar(&cfg.AgentEndPoint, "a", "localhost:8080", "address and port to run agent")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "frequency of report metrics")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "the frequency of the metric survey")
	flag.StringVar(&cfg.AgentKey, "k", "", "the agent key for HashSHA256")
	flag.Int64Var(&cfg.RateLimit, "l", 5, "the maximum number of metrics to report at once")
	flag.Parse()

	var err error

	if envAgentEndPoint := os.Getenv("ADDRESS"); envAgentEndPoint != "" {
		cfg.AgentEndPoint = envAgentEndPoint
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		cfg.ReportInterval, err = strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.PollInterval, err = strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if envAgentKey := os.Getenv("KEY"); envAgentKey != "" {
		cfg.AgentKey = envAgentKey
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		cfg.RateLimit, err = strconv.ParseInt(envRateLimit, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
