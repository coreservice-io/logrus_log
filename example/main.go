package main

import (
	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log"
)

func main() {
	//default is info level
	//return the log interface implemented instance
	// llog, err := logrus_log.NewWithFile("./logs", 1, 20, 30)
	llog, err := logrus_log.New()
	if err != nil {
		panic(err.Error())
	}

	llog.SetLevel(log.TraceLevel)

	llog.Traceln("trace log")
	llog.Debugln("debug log")
	llog.Infoln("info log")
	llog.Warnln("warn log")
	llog.Errorln("error log")
	//log.Fatalln("fatal log")
	//log.Panicln("panic log")

	llog.PrintLastN(100, []log.LogLevel{log.ErrorLevel, log.InfoLevel})
}
