package main

import "flag"

var serverEndPoint string

func parseFlag() {
	flag.StringVar(&serverEndPoint, "a", "localhost:8080", "address and port to run server")

	flag.Parse()
}
