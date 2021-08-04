package ztOrm

import (
	"reflect"
	"strings"
	"testing"
)

func Test_getBySomeTagParams(t *testing.T) {
	// TODO
	if err := getBySomeTagParams(obj, DefaultPostgresEngine, result, "", "xorm", "class_id"); err != nil {
		t.Fatalf("Run %s Err:[%s].", t.Name(), err.Error())
	}
}

func Test_getBySomeTagParams1(t *testing.T) {
	str := "xorm:notnull 'name'"
	first := strings.Index(str, "'")
	last := strings.LastIndex(str, "'")
	first++
	t.Log(first)
	t.Log(last)
	t.Log(str[first:last])
}

func TestClass_GetAllBySomeParams(t *testing.T) {
	result := []string{"1", "2"}

	v := reflect.ValueOf(&result)
	switch v.Kind() {
	case reflect.Ptr:
		result1 := v.Type()
		k := reflect.ValueOf(result1).Kind()
		t.Log(k)
		if k == reflect.Slice {
			t.Log(true)
		} else {
			t.Log(false)
		}
	case reflect.Slice:
		t.Log(true)
	default:
	}
}
