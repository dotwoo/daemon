package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	ds := NewSample()
	setupLog(ds.GetLogFile())
	Run(ds)

}

func setupLog(logFN string) {
	// TODO: 也需要写到 syslogd
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

}
