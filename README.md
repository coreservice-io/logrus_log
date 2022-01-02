# ULog_logrus

##### support both linux mac andwindows

### install
```
go get "github.com/universe-30/ULog_logrus"
```

### example

```go
package main

import (
	"fmt"

	"github.com/universe-30/ULog"
	"github.com/universe-30/ULog_logrus"
)

func main() {
	//default is info level
	ulog_logrus, err := ULog_logrus.New("./logs", 2, 20, 30)
	if err != nil {
		panic(err.Error())
	}

	//ulog_logrus implements the ULog interface
	var ulog ULog.Logger
	ulog = ulog_logrus

	ulog.SetLevel(ULog.TraceLevel)
	ulog.Traceln("trace log")
	ulog.Debugln("debug log")
	ulog.Infoln("info log")
	ulog.Warnln("warn log")
	ulog.Errorln("error log")
	//ulog.Fatalln("fatal log")
	//ulog.Panicln("panic log")

	//ulog_logrus extended functions
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	//all logs include all types :debug ,info ,warning ,error,panic ,fatal
	ulog_logrus.PrintLastN_AllLogs(100)
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
	//err logs include all types :,error,panic ,fatal
	ulog_logrus.PrintLastN_ErrLogs(100)
	fmt.Println("////////////////////////////////////////////////////////////////////////////////")
}

```