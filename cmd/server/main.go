package main

import (
	f "github.com/shestooy/go-musthave-metrics-tpl.git/internal/flags"
	s "github.com/shestooy/go-musthave-metrics-tpl.git/internal/server/httpserver"
)

func main() {
	err := f.ParseServerFlags()
	if err != nil {
		panic(err)
	}

	err = s.Start()
	if err != nil {
		panic(err)
	}
}
