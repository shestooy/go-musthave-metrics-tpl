package main

import (
	"log"

	a "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent"
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
)

func main() {
	err := f.ParseAgentFlag()
	if err != nil {
		log.Fatal("parse flag for agent failed")
	}

	if err = a.Start(); err != nil {
		log.Fatal(err)
	}
}
