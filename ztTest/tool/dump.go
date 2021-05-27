package tool

import (
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

//const dbFilePath = "../../test.db"

func DumpData(dbFilePath, targetDir string) (err error) {
	engine, err := NewEngine(dbFilePath)
	if err != nil {
		return
	}

	dumper, err := testfixtures.NewDumper(
		testfixtures.DumpDatabase(engine.DB().DB),
		testfixtures.DumpDialect("sqlite"),    // or your database of choice.
		testfixtures.DumpDirectory(targetDir), // api name that wait to test.
	)
	if err != nil {
		return
	}
	if err := dumper.Dump(); err != nil {
		return
	}
	return
}

func NewEngine(sqliteFilePath string) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("sqlite3", sqliteFilePath+"?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(true)
	return
}
