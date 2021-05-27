package main

import (
	"fmt"
	"github.com/zput/zxcTool/ztTest"
	"io/ioutil"
	"xorm.io/xorm"
)

type Example struct {
	Id   int64  `xorm:"pk autoincr bigint" json:"id"`
	Name string `xorm:"varchar(10) not null" json:"name"`
}

func main() {
	var (
		path = "get_example_by_id"
		//path        = "/Users/edz/CODE/Self/zxcTool/ztTest/example/get_example_by_id"
		tablePrefix = "table_"
	)

	_, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if engine, err := example1(path, tablePrefix); err != nil {
		panic(err)
	} else {
		var name string
		if _, err := engine.SQL("select name from table_example").Get(&name); err != nil {
			panic(err)
		}
		if name != "ztTest" {
			panic(fmt.Sprintf("expect ztTest, but get %s", name))
		} else {
			fmt.Println("(- v -), pass")
		}
	}

}

func example1(path string, tablePrefix string) (engine *xorm.Engine, err error) {
	var f ztTest.IFixtureServe
	f, err = ztTest.New(ztTest.SetFixturePath(path),
		ztTest.SetFixtureTablePrefix(tablePrefix))
	if err != nil {
		return
	}

	if err = f.Sync(
		new(Example),
	); err != nil {
		return
	}

	if err = f.Load(); err != nil {
		return
	}

	engine = f.Engine()
	return
}
