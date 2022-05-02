# logrus_log

##### support both linux mac andwindows

### install
```
go get "github.com/universe-30/logrus_log"
```

### example

```go
package main

import (
	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log"
)

func main() {
	//default is info level
	llog, err := logrus_log.New("./logs", 1, 20, 30)
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

```