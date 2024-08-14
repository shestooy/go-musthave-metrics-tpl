package main

import (
	a "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"log"
)

func main() {
	f.ParseAgentFlag()
	err := a.Start()
	if err != nil {
		log.Fatal("send metrics failed")
	}
}
