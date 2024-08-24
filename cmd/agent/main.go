package main

import (
	a "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	"log"
)

func main() {
	err := f.ParseAgentFlag()
	if err != nil {
		log.Fatal("parse flag for agent failed")
	}

	a.Start()
}
