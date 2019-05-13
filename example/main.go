package main

import (
	"github.com/dotwoo/daemon"
	log "github.com/sirupsen/logrus"
)

func main() {
	ds := NewSample()
	setupLog(ds.GetLogFile())
	daemon.Run(ds)

}

func setupLog(logFN string) {
	// TODO: 也需要写到 syslogd
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

}
