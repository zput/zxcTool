package ztOrm

import (
	"github.com/astaxie/beego/config"
	"testing"
)

func TestNewDbEngine(t *testing.T) {
	cfg, err := config.NewConfig("ini", "../../conf/server.ini")
	if err != nil {
		t.Fatalf("%s.Err:%v. ", t.Name(), err)
	}
	engine := NewDbEngine(cfg)
	if _, err := engine.GetPostgresEngine(); err != nil {
		t.Fatalf("%s.Err:%v. ", t.Name(), err)
	}

	if _, err := engine.GetRedisEngine(); err != nil {
		t.Fatalf("%s.Err:%v. ", t.Name(), err)
	}
}
