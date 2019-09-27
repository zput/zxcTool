package zt_formatter

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"runtime"
)

type CallBackFoo func (*runtime.Frame)(function string, file string)
type WriteEntry func (entry *logrus.Entry, callerPrettyfier CallBackFoo)

type FormaterOperatorInterface interface {
	WriteCommonInfo(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry)
	WriteField(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry)
	WriteMessages(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry)
}

type FormatterZTInterface interface {
	WriteEntry(entry *logrus.Entry)(fileVal, funcVal string)
	WriteFields(b *bytes.Buffer, entry *logrus.Entry)
	WriteOrderedFields(b *bytes.Buffer, entry *logrus.Entry)
    WriteField(b *bytes.Buffer, entry *logrus.Entry, field string)

	// no need to rewrite
	GetTimestampFormat()(string)
	GetNoColors()(bool)
	GetShowFullLevel()(bool)
	GetNoFieldsColors()(bool)
	GetFieldsOrder()([]string)
	GetTrimMessages()(bool)
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.WarnLevel:
		return colorBlue
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorGray
	}
}

func getFileFooColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
