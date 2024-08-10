package main

import (
	"flag"
	"os"
)

var serverEndPoint string

func parseFlag() {
	flag.StringVar(&serverEndPoint, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envServerEndPoint := os.Getenv("ADDRESS"); envServerEndPoint != "" {
		serverEndPoint = envServerEndPoint
	}
}
