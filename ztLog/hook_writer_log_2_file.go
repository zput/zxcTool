package ztLog

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// WriterToFileHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterToFileHook struct {
	LogNamePrefix string
	Writer        io.Writer
	LogLevels     []log.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterToFileHook) Fire(entry *log.Entry) error {
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
		hook.Writer = os.Stderr
		hook.Writer.Write([]byte("Failed to log to file, using default stderr"))
	}
	defer file.Close()

	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterToFileHook) Levels() []log.Level {
	return hook.LogLevels
}
