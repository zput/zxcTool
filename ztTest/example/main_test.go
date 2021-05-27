package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestExample(t *testing.T) {
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
			fmt.Println("- v -, pass")
		}
	}

}
