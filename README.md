
# Powerful-logrus-formatter(input error line and save log to files)


Human-readable log formatter:

![Screenshot](https://github.com/zput/myPicLib/raw/master/GO/github/log.demo.png)


## Charactor

1. print the name and line about file and the lastest function name while this log
2. save log to file


## Configuration:

```go
type ZtFormatter struct{
  nested.Formatter
    // CallerPrettyfier can be set by the user to modify the content
    // of the function and file keys in the json data when ReportCaller is
    // activated. If any of the returned value is the empty string the
    // corresponding key will be removed from json fields.
  CallerPrettyfier CallBackFoo

  FormaterOperator FormaterOperatorInterface
}

```

## Example

```go
package main

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"path"
	"runtime"
)

func main() {
	var exampleFormatter = &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			//HideKeys: true,
			FieldsOrder: []string{"component", "category"},
		},
	}
	printDemo(exampleFormatter, "hello world")
}

func printDemo(f logrus.Formatter, title string) {
	l := logrus.New()

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

```

### use zxclog

```
	err = ztUtil.CodedNoPtrErrorf(errno.EN_InnerServer, errno.EN_InnerServer, "%v", err)
	logs.Error(err)
	return



	if err != nil {
		logs.Error("commonInfo:%s; PostSaleCenterUserActivity; error:%s", commonInfo, err.Error())
		if errReal, ok := err.(*ztUtil.Error); ok {
			util.OnError(this.Ctx, *errReal.Code, errReal.Message)
			return
		} else {
			util.OnError(this.Ctx, errno.EN_InnerServer, errno.EM_InnerServer)
			return
		}
	}
```


> more example please reference: /ztLog/example



## Reference 

thanks [nested-logrus-formatter](https://github.com/antonfisher/nested-logrus-formatter)

## Welcome PR


<!---
[![Build Status](https://travis-ci.org/antonfisher/nested-logrus-formatter.svg?branch=master)](https://travis-ci.org/antonfisher/nested-logrus-formatter)
[![Go Report Card](https://goreportcard.com/badge/github.com/antonfisher/nested-logrus-formatter)](https://goreportcard.com/report/github.com/antonfisher/nested-logrus-formatter)
[![GoDoc](https://godoc.org/github.com/antonfisher/nested-logrus-formatter?status.svg)](https://godoc.org/github.com/antonfisher/nested-logrus-formatter)
-->

