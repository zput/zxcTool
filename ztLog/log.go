package ztLog

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/astaxie/beego/config"
	log "github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var (
	GlobalInterval int
)

func LogInitialize(cfg config.Configer) {
	logNamePrefix := cfg.String("log::log_name_prefix")
	logFormatter := cfg.DefaultString("log::log_formatter", "NESTEDFormatter")
	logLevel := cfg.DefaultInt("log::log_level", 7)
	whetherWriteToFile := cfg.DefaultBool("log::log_whether_write_to_file", false)
	logFilePath := cfg.DefaultString("log::log_file_path", "../temp_log/")
	SetupSTDLogs(whetherWriteToFile, logFilePath+logNamePrefix, logFormatter, logLevel)
}

// SetupSTDLogs adds hooks to send logs to different destinations depending on level
func SetupSTDLogs(whetherWriteToFile bool, logNamePrefix, logFormatter string, logLevel int) {
	var choicedFormatter log.Formatter

	switch logFormatter {
	case "NESTEDFormatter":
		choicedFormatter = &zt_formatter.ZtFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			Formatter: nested.Formatter{
				//HideKeys: true,
				FieldsOrder: []string{"component", "category"},
			},
		}
	case "JSONFormatter":
		choicedFormatter = &log.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
	}

	SetupLogsCanExpand(log.StandardLogger(), choicedFormatter, whetherWriteToFile, logNamePrefix, logLevel)
}

func SetupLogsCanExpand(lPtr *log.Logger, f log.Formatter, whetherWriteToFile bool, logNamePrefix string, logLevel int) {
	lPtr.SetLevel(log.Level(logLevel - 1))
	lPtr.SetReportCaller(true)
	lPtr.Formatter = f
	lPtr.SetOutput(ioutil.Discard)  // Send all logs to nowhere by default
	lPtr.AddHook(&WriterToFileHook{ // Send logs with level higher than warning to stderr
		WhetherWriteToFile: whetherWriteToFile,
		LogNamePrefix:      logNamePrefix,
		Writer:             os.Stdout,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})
}

//func SetupLogs(logNamePrefix, logFormatter string, logLevel int) {
//	/*	err := logrus_mate.Hijack(
//			log.StandardLogger(),
//			logrus_mate.ConfigString(
//				`{formatter.name = "json"}`,
//			),
//		)
//		if err != nil{
//			panic(fmt.Sprintf("err:%s", err.Error()))
//		}*/
//
//	log.SetLevel(log.Level(logLevel - 1))
//
//	log.SetReportCaller(true)
//
//	switch logFormatter {
//	case "NESTEDFormatter":
//		log.StandardLogger().Formatter = &zt_formatter.ZtFormatter{
//			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
//				filename := path.Base(f.File)
//				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
//			},
//			Formatter: nested.Formatter{
//				//HideKeys: true,
//				FieldsOrder: []string{"component", "category"},
//			},
//		}
//	case "JSONFormatter":
//		log.StandardLogger().Formatter = &log.JSONFormatter{
//			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
//				filename := path.Base(f.File)
//				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
//			},
//		}
//	}
//
//	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default
//
//	log.AddHook(&WriterToFileHook{ // Send logs with level higher than warning to stderr
//		LogNamePrefix: logNamePrefix,
//		Writer:        os.Stdout,
//		LogLevels: []log.Level{
//			log.PanicLevel,
//			log.FatalLevel,
//			log.ErrorLevel,
//			log.WarnLevel,
//			log.InfoLevel,
//			log.DebugLevel,
//			log.TraceLevel,
//		},
//	})
//
//}

/*
func main(){
	SetupLogs()
    for i:=0; i<1000; i++{
		log.Error(i, "---xxx\n")
    	time.Sleep(1*time.Second)
	}
}
*/
