# Logrus

##### support both linux mac andwindows

### install
```
go get "github.com/universe-30/logrusULog"
```

### example

```go
package main
 
import (
	"github.com/coreservice-io/LogrusULog"
	"github.com/coreservice-io/ULog"
)

func main() {
	//default is info level
	ulog, err := LogrusULog.New("./logs", 1, 20, 30)
	if err != nil {
		panic(err.Error())
	}

	ulog.SetLevel(ULog.TraceLevel)

	ulog.Traceln("trace log")
	ulog.Debugln("debug log")
	ulog.Infoln("info log")
	ulog.Warnln("warn log")
	ulog.Errorln("error log")
	//ulog.Fatalln("fatal log")
	//ulog.Panicln("panic log")

	ulog.PrintLastN(100, []ULog.LogLevel{ULog.ErrorLevel, ULog.InfoLevel})
}

```