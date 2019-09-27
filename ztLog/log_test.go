package ztLog_test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog"
	"testing"
)

func TestLog(t *testing.T){
	ztLog.SetupLogs("./xxx/xxxx", "NESTEDFormatter", 7)

	log.Info("hello")

	fmt.Println("--->>>")
	t.Error("xxx")
}
