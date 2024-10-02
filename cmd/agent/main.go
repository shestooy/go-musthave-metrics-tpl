package main

import (
	"log"

	a "github.com/shestooy/go-musthave-metrics-tpl.git/internal/agent"
)

func main() {
	if err := a.Start(); err != nil {
		log.Fatal(err)
	}
}
