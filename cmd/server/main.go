package main

import (
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	s "github.com/shestooy/go-musthave-metrics-tpl.git/internal/httpserver"
)

func main() {
	f.ParseServerFlags()

	err := s.Start()
	if err != nil {
		panic(err)
	}
}
