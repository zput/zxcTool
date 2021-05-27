package ztTest

import (
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

func DropTables(engine *xorm.Engine, tables ...interface{}) (err error) {
	for _, v := range tables {
		err = engine.DropTables(v)
		if err != nil {
			return
		}
	}
	return
}

func SyncTables(engine *xorm.Engine, tables ...interface{}) (err error) {
	for _, v := range tables {
		err = engine.Sync2(v)
		if err != nil {
			return
		}
	}
	return
}

type fixture struct {
	*testfixtures.Loader
	*xorm.Engine
}

func newEngine(tablePrefix string) (res *xorm.Engine, err error) {
	res, err = xorm.NewEngine("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return
	}
	res.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, tablePrefix))
	res.ShowSQL(true)
	res.SetLogLevel(xlog.LOG_DEBUG)
	return
}

func newFixture(path, tablePrefix string) (*fixture, error) {

	engine, err := newEngine(tablePrefix)
	if err != nil {
		return nil, err
	}

	var para = []func(*testfixtures.Loader) error{
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Database(engine.DB().DB), // You database connection
		testfixtures.Dialect("sqlite"),        // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		//testfixtures.Directory(path),          // The directory containing the YAML files
	}
	if len(path) > 0 {
		para = append(para, testfixtures.Directory(path))
	}

	f, err := testfixtures.New(para...)
	if err != nil {
		return nil, err
	}
	return &fixture{f, engine}, nil
}
