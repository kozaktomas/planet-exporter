package main

import (
	"flag"
	"log"
)

var pPort int // pPort is a port for Prometheus exporter

func init() {
	flag.IntVar(&pPort, "port", 9042, "Exporter port to listen on")
}

func main() {
	flag.Parse()

	monitor := NewMonitor()
	solarSystem, err := NewSolarSystem(monitor)
	if err != nil {
		log.Fatalf("Failed to load planet data: %v", err)
	}
	router := NewRouter(solarSystem, monitor)
	StartServer(router, pPort)
}
