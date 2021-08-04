package ztUtil

import (
	"testing"
)

type Class struct {
	Id                 int64  `xorm:"pk autoincr 'id'"`
	Tag                string `xorm:"notnull 'tag'"`
	Len                int    `xorm:"notnull default(0) 'len'"`
	Creator            int64  `xorm:"notnull 'creator'"`
	Capacity           int    `xorm:"notnull default(0) 'capacity'"`
	SeriesID           int64  `xorm:"notnull default(0) 'series_id'"`
	BeginTime          int64  `xorm:"notnull 'begin_time'"`
	CreateTime         int64  `xorm:"created"`
	UpdateTime         int64  `xorm:"updated"`
	RemainingClassNum  int    `xorm:"notnull default(1) 'remaining_class'"`
	Consumption        int    `xorm:"notnull default(1) 'consumption'"`
	CourseSerialNumber int    `xorm:"notnull default(1) 'course_serial_number'"`
}

func (c *Class) TableName() string {
	return "class"
}

func TestGetStatusTagName(t *testing.T) {
	obj := &Class{}
	t.Log(GetTagName(obj, "Tag", "xorm"))
}

func TestSetFieldValue(t *testing.T) {
	obj := &Class{}
	var id int64
	id = 1

	flag := SetFieldValue("Id", obj, id)
	t.Log(flag)
	t.Logf("Class:%+v", obj)
}
