package ztTest

import (
	"fmt"
	"io/ioutil"
	"xorm.io/xorm"
)

/*
s := ztTest.New() // ---configuration---
s.Construct()
s.Load()
s.Engine()
*/

func New(o ...FixtureOption) (IFixtureServe, error) {
	return NewFixture(o...)
}

type IFixtureServe interface {
	Sync(tables ...interface{}) error
	Load() error
	Engine() *xorm.Engine
}

// NewFixture example: NewFixture(path).Prepare()
// self fixture struct
func NewFixture(o ...FixtureOption) (res *Fixture, err error) {
	ret := new(Fixture)
	for i := range o {
		if err = o[i](ret); err != nil {
			return
		}
	}

	jiuWMySql, err := newFixture(ret.path, ret.tablePrefix)
	if err != nil {
		return
	}
	ret.jiuWMySql = jiuWMySql
	return ret, nil
}

type FixtureOption func(o *Fixture) error

// SetFixturePath set path in order to load YAML files from a given directory.
func SetFixturePath(path string) FixtureOption {
	return func(o *Fixture) (err error) {
		_, err = ioutil.ReadDir(path)
		if err != nil {
			err = fmt.Errorf(`ztTest: could not stat directory "%s": %w`, path, err)
			return
		}
		o.path = path
		return
	}
}

func SetFixtureTablePrefix(tablePrefix string) FixtureOption {
	return func(o *Fixture) (err error) {
		o.tablePrefix = tablePrefix
		return
	}
}

type Fixture struct {
	path        string // path in order to load YAML files from a given directory.
	tablePrefix string
	jiuWMySql   *fixture
}

func (f *Fixture) Sync(tables ...interface{}) error {
	return SyncTables(f.jiuWMySql.Engine, // 表初始化到数据库
		tables...)
}

func (f *Fixture) Load() error {
	if err := f.jiuWMySql.Load(); err != nil {
		return fmt.Errorf("cannot load jiuWMySql fixtures, err: %+v", err)
	}
	return nil
}

func (f *Fixture) Engine() *xorm.Engine {
	return f.jiuWMySql.Engine
}
