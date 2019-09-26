package ztLog

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	log "github.com/sirupsen/logrus"
)

var (
	GlobalInterval int
)

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	LogNamePrefix string
	Writer        io.Writer
	LogLevels     []log.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	nowInterval := time.Now().Day()
	if nowInterval != GlobalInterval {
		GlobalInterval = nowInterval
	}

	fileName := fmt.Sprintf("%s_%s_%d.log", hook.LogNamePrefix, time.Now().Format("2006-01"), GlobalInterval)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		hook.Writer = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	defer file.Close()

	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []log.Level {
	return hook.LogLevels
}

// setupLogs adds hooks to send logs to different destinations depending on level
func SetupLogs(logNamePrefix, logFormatter string, logLevel int) {
	/*	err := logrus_mate.Hijack(
			log.StandardLogger(),
			logrus_mate.ConfigString(
				`{formatter.name = "json"}`,
			),
		)
		if err != nil{
			panic(fmt.Sprintf("err:%s", err.Error()))
		}*/

	log.SetLevel(log.Level(logLevel - 1))

	log.SetReportCaller(true)

	switch logFormatter {
	case "NESTEDFormatter":
		log.StandardLogger().Formatter = &zt_formatter.ZtFormatter{
			Formatter: nested.Formatter{
				HideKeys: true,
				FieldsOrder: []string{"component", "category"},
			},
		}
	case "JSONFormatter":
		log.StandardLogger().Formatter = &log.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
	}

	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	log.AddHook(&WriterHook{ // Send logs with level higher than warning to stderr
		LogNamePrefix: logNamePrefix,
		Writer:        os.Stdout,
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

/*
func main(){
	SetupLogs()
    for i:=0; i<1000; i++{
		log.Error(i, "---xxx\n")
    	time.Sleep(1*time.Second)
	}
}
*/

func LogInitialize(cfg config.Configer) {
	logNamePrefix := cfg.String("log::log_name_prefix")
	logFormatter := cfg.DefaultString("log::log_formatter", "NESTEDFormatter")
	logLevel := cfg.DefaultInt("log::log_level", 7)
	SetupLogs("../temp_log/"+logNamePrefix, logFormatter, logLevel)
}
