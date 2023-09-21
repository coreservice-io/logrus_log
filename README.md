# logrus_log

##### logrus_log implement the log interface

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
	//return the log interface implemented instance

	// if need to log to file
	llog, err := logrus_log.NewWithFile("./logs", 1, 20, 30)
	// if do not want to log to file, only log to console
	//  llog, err := logrus_log.New()
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