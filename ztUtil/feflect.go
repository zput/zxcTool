package ztUtil

import (
	"github.com/go-xorm/xorm"
	"reflect"
	"strings"
)

func GetTagName(obj xorm.TableName, name, tag string) string {
	subStr := "'"
	s := reflect.TypeOf(obj).Elem()
	nameField, isExists := s.FieldByName(name)
	if !isExists {
		return ""
	}
	valueXormTagName := nameField.Tag.Get(tag)

	if strings.Contains(valueXormTagName, subStr) {
		first := strings.Index(valueXormTagName, subStr)
		first++
		last := strings.LastIndex(valueXormTagName, subStr)
		valueXormTagName = valueXormTagName[first:last]
	}
	return valueXormTagName
}

func SetFieldValue(fieldName string, obj, value interface{}) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return false
	}

	if canSet := v.Elem().FieldByName(fieldName).CanSet(); !canSet {
		return canSet
	}

	vValue := reflect.ValueOf(value)
	v.Elem().FieldByName(fieldName).Set(vValue)
	return true
}
