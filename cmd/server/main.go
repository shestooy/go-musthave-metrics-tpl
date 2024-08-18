package main

import (
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	s "github.com/shestooy/go-musthave-metrics-tpl.git/internal/httpserver"
)

func main() {
	flags.ParseServerFlags()

	err := s.Start()
	if err != nil {
		panic(err)
	}
}
