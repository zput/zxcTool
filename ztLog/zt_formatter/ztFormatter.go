package zt_formatter

import (
	"bytes"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"runtime"
	"sort"
	"strings"
	"time"
)
type CallBackFoo func (*runtime.Frame)(function string, file string)

type WriteEntry func (entry *logrus.Entry, callerPrettyfier CallBackFoo)

type ZtFormatter struct{
	nested.Formatter

    // CallerPrettyfier can be set by the user to modify the content
    // of the function and file keys in the json data when ReportCaller is
    // activated. If any of the returned value is the empty string the
    // corresponding key will be removed from json fields.
    CallerPrettyfier CallBackFoo
}


// rewrite Format an log entry by zxc
func (f *ZtFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	level := strings.ToUpper(entry.Level.String())

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("] ")

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// add for zxc --------------- ---
	f.writeEntry(entry)

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}


func (f *ZtFormatter)writeEntry(entry *logrus.Entry) {

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			entry.Data[logrus.FieldKeyFunc] = funcVal
			var tempSliceString = make([]string, 1)
			tempSliceString = append(tempSliceString, logrus.FieldKeyFunc)
			tempSliceString = append(tempSliceString, f.FieldsOrder...)
			f.FieldsOrder = tempSliceString
		}
		if fileVal != "" {
			entry.Data[logrus.FieldKeyFile] = fileVal
			var tempSliceString = make([]string, 1)
			tempSliceString = append(tempSliceString, logrus.FieldKeyFile)
			tempSliceString = append(tempSliceString, f.FieldsOrder...)
			f.FieldsOrder = tempSliceString
		}
	}

}


func (f *ZtFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *ZtFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *ZtFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}




















