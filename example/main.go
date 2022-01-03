package main

import (
	"github.com/universe-30/Logrus"
	"github.com/universe-30/ULog"
)

func main() {
	//default is info level
	log_logrus, err := Logrus.New("./logs", 1, 20, 30)
	if err != nil {
		panic(err.Error())
	}

	//ulog_logrus implements the ULog interface
	var ulog ULog.Logger
	ulog = log_logrus

	ulog.SetLevel(ULog.TraceLevel)

	ulog.Traceln("trace log")
	ulog.Debugln("debug log")
	ulog.Infoln("info log")
	ulog.Warnln("warn log")
	ulog.Errorln("error log")
	//ulog.Fatalln("fatal log")
	//ulog.Panicln("panic log")

	//ulog_logrus extended functions
	//all logs include all types :debug ,info ,warning ,error,panic ,fatal
	log_logrus.PrintLastN_AllLogs(100)
	//err logs include all types :,error,panic ,fatal
	log_logrus.PrintLastN_ErrLogs(100)
}
