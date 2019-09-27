package zt_formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DefaultFormaterrOperator struct{
}

func (*DefaultFormaterrOperator)WriteCommonInfo(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry){

	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.GetTimestampFormat()
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// add for zxc --------------- ---
	fileVal, funcVal := f.WriteEntry(entry)
	if fileFooColor := getFileFooColorByLevel(entry.Level); fileFooColor == colorRed{
		bPtr.WriteString(fmt.Sprintf("\x1b[%dm[%s::%s]\x1b[0m ", fileFooColor, fileVal, funcVal))
	}else{
		bPtr.WriteString(fmt.Sprintf("[%s::%s] ", fileVal, funcVal))
	}

	// write time
	bPtr.WriteString(entry.Time.Format(timestampFormat))

	// write level
	level := strings.ToUpper(entry.Level.String())

	if !f.GetNoColors() {
		fmt.Fprintf(bPtr, "\x1b[%dm", levelColor)
	}

	bPtr.WriteString(" [")
	if f.GetShowFullLevel() {
		bPtr.WriteString(level)
	} else {
		bPtr.WriteString(level[:1])
	}
	bPtr.WriteString("] ")

	if !f.GetNoColors() && f.GetNoFieldsColors() {
		bPtr.WriteString("\x1b[0m")
	}
}


func (*DefaultFormaterrOperator)WriteField(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry){

	// write fields
	if f.GetFieldsOrder() == nil {
		f.WriteFields(bPtr, entry)
	} else {
		f.WriteOrderedFields(bPtr, entry)
	}

	if !f.GetNoColors() && !f.GetNoFieldsColors() {
		bPtr.WriteString("\x1b[0m")
	}

}

func (*DefaultFormaterrOperator)WriteMessages(f FormatterZTInterface, bPtr *bytes.Buffer, entry *logrus.Entry){
	// write message
	if f.GetTrimMessages() {
		bPtr.WriteString(strings.TrimSpace(entry.Message))
	} else {
		bPtr.WriteString(entry.Message)
	}
}
