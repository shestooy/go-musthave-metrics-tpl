package main

import (
	s "github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/httpserver"
)

func main() {
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
