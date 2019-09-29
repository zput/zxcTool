package zt_formatter

import (
	"bytes"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"sort"
)

type ZtFormatter struct{

	nested.Formatter
    // CallerPrettyfier can be set by the user to modify the content
    // of the function and file keys in the json data when ReportCaller is
    // activated. If any of the returned value is the empty string the
    // corresponding key will be removed from json fields.
    CallerPrettyfier CallBackFoo

	FormaterOperator FormaterOperatorInterface
}

func(f *ZtFormatter)GetTimestampFormat()(string){
	return f.TimestampFormat
}

func(f *ZtFormatter)GetNoColors()(bool){
	return f.NoColors
}

func(f *ZtFormatter)GetShowFullLevel()(bool){
	return f.ShowFullLevel
}

func(f *ZtFormatter)GetNoFieldsColors()(bool){
	return f.NoFieldsColors
}

func(f *ZtFormatter)GetFieldsOrder()([]string){
	return f.FieldsOrder
}


func(f *ZtFormatter)GetTrimMessages()(bool){
	return f.TrimMessages
}



// rewrite Format an log entry by zxc
func (f *ZtFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// output buffer
	b := &bytes.Buffer{}

	if f.FormaterOperator == nil{
		f.FormaterOperator = new(DefaultFormaterrOperator)
	}

	f.FormaterOperator.WriteCommonInfo(f, b, entry)
	f.FormaterOperator.WriteField(f, b, entry)
	f.FormaterOperator.WriteMessages(f, b, entry)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *ZtFormatter) WriteEntry(entry *logrus.Entry)(fileVal, funcVal string){
	if entry.HasCaller() {
		funcVal = entry.Caller.Function
		fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		// fmt.Println(funcVal, fileVal)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
	}
	return
}

func (f *ZtFormatter) WriteFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.WriteField(b, entry, field)
		}
	}
}

func (f *ZtFormatter) WriteOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.WriteField(b, entry, field)
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
			f.WriteField(b, entry, field)
		}
	}
}

func (f *ZtFormatter) WriteField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}
