package main

import (
	"flag"
)

var pPort int // pPort is a port for Prometheus exporter

func init() {
	flag.IntVar(&pPort, "port", 9042, "Exporter port to listen on")
}

func main() {
	flag.Parse()

	monitor := NewMonitor()
	router := NewRouter(monitor)
	StartServer(router, pPort)
}
