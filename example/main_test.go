package example

import (
	"fmt"
	"testing"

	"github.com/universe-30/ULog_logrus"
	"github.com/universe-30/UUtils/path_util"
)

func Test_main(t *testing.T) {
	//default is info level
	ulog, err := ULog_logrus.New(path_util.GetAbsPath("logs"), 2, 20, 30)
	if err != nil {
		panic(err.Error())
	}

	ulog.WithFields(ULog_logrus.Fields{
		"f1": "1",
		"f2": "2",
	}).Errorf("Total xxx Error Fileds : %d", 2)

	ulog.WithFields(ULog_logrus.Fields{
		"f1": "1",
		"f2": "2",
	}).Warnf("Total  yy Warn Fileds : %d", 2)

	ulog.SetLevel(ULog_logrus.DebugLevel)

	ulog.WithFields(ULog_logrus.Fields{
		"f1": "1",
		"f2": "2",
	}).Infof("Total zzz Fileds : %d", 2)

	ulog.WithFields(ULog_logrus.Fields{
		"f1": "1",
		"f2": "2",
	}).Debugf("Total Debug Fileds : %d", 2)

	ulog.SetLevel(ULog_logrus.TraceLevel)

	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	//all logs include all types :debug ,info ,warning ,error,panic ,fatal
	ulog.PrintLastN_AllLogs(100)
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	//err logs include all types :,error,panic ,fatal
	ulog.PrintLastN_ErrLogs(100)
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")

	ulog.Warnln("this is warn ln")
	ulog.Warnf("this is warnf %d", 123)

	ulog.Traceln("this is warn ln")
}
