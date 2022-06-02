package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Zedronar/go-dummy-app.git/external/metrics"
	"github.com/Zedronar/go-dummy-app.git/internal"
)

const (
	metricNameServiceStarts = "service.starts"
)

func main() {
	internal.InitLiveServer()

	metrics.Increment(metricNameServiceStarts)

	runService()
}

func runService() {
	osStopChannel := make(chan os.Signal)
	signal.Notify(osStopChannel, os.Interrupt)
	log.Println("service started")
	<-osStopChannel
	log.Println("service stopping...")
}
