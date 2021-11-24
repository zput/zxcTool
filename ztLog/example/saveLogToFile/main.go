package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog"
)

func main() {
	logNamePrefix := "log_name_prefix"
	logFormatter := "NESTEDFormatter"
	logLevel := 7
	ztLog.SetupSTDLogs(true, "./"+logNamePrefix, logFormatter, logLevel)

	printDemo(logrus.StandardLogger(), nil, "hello world")
}

func printDemo(l *logrus.Logger, f logrus.Formatter, title string) {
	//l := logrus.New()

	l.SetLevel(logrus.DebugLevel)
	l.SetReportCaller(true)

	if f != nil {
		l.SetFormatter(f)
	}

	l.Infof("this is %v demo", title)

	lWebServer := l.WithField("component", "web-server")
	lWebServer.Info("starting...")

	lWebServerReq := lWebServer.WithFields(logrus.Fields{
		"req":   "GET /api/stats",
		"reqId": "#1",
	})

	lWebServerReq.Info("params: startYear=2048")
	lWebServerReq.Error("response: 400 Bad Request")

	lDbConnector := l.WithField("category", "db-connector")
	lDbConnector.Info("connecting to db on 10.10.10.13...")
	lDbConnector.Warn("connection took 10s")

	l.Info("demo end.")
}
