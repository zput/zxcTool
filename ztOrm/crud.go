package ztOrm

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	util "github.com/zput/zxcTool/ztUtil"
	"gopkg.in/mgo.v2"
	"reflect"
	"strings"
)

var (
	DefaultRedis          *redis.Client
	DefaultPostgresEngine *xorm.Engine
	DefaultMgoSess        *mgo.Session
	mgoDB                 string
)

var DURATION_AFTER_CLASS = "INTERVAL '240 min'"

func insert(tableName string, engine *xorm.Engine, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		obj = &obj
	}

	if engine == nil {
		return errors.New("Nil engine. ")
	}
	var id int64
	util.SetFieldValue("Id", obj, id) //将ID设置为0
	affected, err := engine.Table(tableName).Insert(obj)
	if err != nil {
		return err
	}
	if affected == 0 {
		return DaoErrFailedToInsert
	}
	return nil
}

func insertByTransaction(tableName string, engine *xorm.Session, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		obj = &obj
	}

	if engine == nil {
		return errors.New("Nil engine. ")
	}

	var id int64
	util.SetFieldValue("Id", obj, id) //将ID设置为0
	affected, err := engine.Table(tableName).Insert(obj)
	if err != nil {
		return err
	}
	if affected == 0 {
		return DaoErrFailedToInsert
	}
	return nil
}

func insertAllByTransaction(obj xorm.TableName, engine *xorm.Session, allObj interface{}) error {
	if reflect.ValueOf(allObj).Kind() != reflect.Ptr {
		allObj = &allObj
	}

	affected, err := engine.Table(obj.TableName()).Insert(allObj)
	if err != nil {
		return err
	}

	if affected == 0 {
		return DaoErrFailedToInsert
	}

	return nil
}

func update(obj xorm.TableName, engine *xorm.Engine, ID int64, cols ...string) (int64, error) {
	if engine == nil {
		return 0, errors.New("Nil engine. ")
	}

	if ID == 0 {
		return 0, DaoErrNotId
	}

	return engine.Table(obj.TableName()).Cols(cols...).ID(ID).Update(obj)
}

func updateByTransaction(obj xorm.TableName, session *xorm.Session, ID int64, cols ...string) (int64, error) {
	if ID == 0 {
		return 0, DaoErrNotId
	}
	return session.Table(obj.TableName()).ID(ID).Cols(cols...).Update(obj)
}

func daoDelete(obj xorm.TableName, engine *xorm.Engine, ID int64) (int64, error) {
	if ID == 0 {
		return 0, DaoErrNotId
	}

	if engine == nil {
		return 0, errors.New("Nil engine. ")
	}

	affected, err := engine.Table(obj.TableName()).ID(ID).Delete(obj)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func daoDeleteByTransaction(obj xorm.TableName, session *xorm.Session, ID int64) (int64, error) {
	if ID == 0 {
		return 0, DaoErrNotId
	}

	affected, err := session.Table(obj.TableName()).Delete(obj)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func get(obj xorm.TableName, engine *xorm.Engine) error {
	isExists, err := engine.Table(obj.TableName()).Get(obj)
	if err != nil {
		return err
	}
	if !isExists {
		return DaoErrNotExist
	}
	return nil
}

/**
  tag: eg:json\xorm...
  params: eg: name string 'xorm:"'name'"'
  TODO 其实obj中的数据被赋值就会被xorm自动映射成查询参数
*/
func getBySomeTagParams(obj xorm.TableName, engine *xorm.Engine, result interface{}, orderBy, tag string, params ...string) error {
	subStr := "'"
	format := "%s = ?"
	v := reflect.Indirect(reflect.ValueOf(obj))
	t := v.Type()

	sess := engine.Table(obj.TableName())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		valueXormTagName := field.Tag.Get(tag)
		if strings.Contains(valueXormTagName, subStr) {
			first := strings.Index(valueXormTagName, subStr)
			first++
			last := strings.LastIndex(valueXormTagName, subStr)
			valueXormTagName = valueXormTagName[first:last]
		}

		if util.Contains(valueXormTagName, params...) {
			sess.And(fmt.Sprintf(format, valueXormTagName), value)
		}
	}

	if orderBy != "" {
		sess.OrderBy(orderBy)
	}

	v = reflect.Indirect(reflect.ValueOf(result))
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		return sess.Find(result)
	}

	isExists, err := sess.Get(result)
	if err != nil {
		return err
	}

	if !isExists {
		return DaoErrNotExist
	}

	return nil
}

func getBySomeTagParamsAndTransaction(obj xorm.TableName, session *xorm.Session, result interface{}, orderBy, tag string, params ...string) error {
	subStr := "'"
	format := "%s = ?"
	v := reflect.Indirect(reflect.ValueOf(obj))
	t := v.Type()

	sess := session.Table(obj.TableName())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		valueXormTagName := field.Tag.Get(tag)
		if strings.Contains(valueXormTagName, subStr) {
			first := strings.Index(valueXormTagName, subStr)
			first++
			last := strings.LastIndex(valueXormTagName, subStr)
			valueXormTagName = valueXormTagName[first:last]
		}

		if util.Contains(valueXormTagName, params...) {
			sess.And(fmt.Sprintf(format, valueXormTagName), value)
		}
	}

	if orderBy != "" {
		sess.OrderBy(orderBy)
	}

	v = reflect.Indirect(reflect.ValueOf(result))
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		return sess.Find(result)
	}

	isExists, err := sess.Get(result)
	if err != nil {
		return err
	}

	if !isExists {
		return DaoErrNotExist
	}

	return nil
}
