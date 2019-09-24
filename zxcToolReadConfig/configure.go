package zxcToolReadConfig

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/utils"
	"os"
	"path/filepath"
)

var (
	InnerConfig config.Configer
)

//func init(){
//	ConfiguresInit("/online-compiler/goLanguage/configures/")
//}

// TestBeegoInit is for test package init
func ConfiguresInit(ap string) {
	path := filepath.Join(ap, "conf", "app.conf")
	os.Chdir(ap)

	if err := LoadAppConfig("ini", path); err != nil {
		panic(err)
	}
}

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	if !utils.FileExists(absConfigPath) {
		return fmt.Errorf("the target config file: %s don't exist", configPath)
	}

	InnerConfig, err = config.NewConfig(adapterName, configPath)
	if err != nil {
		return err
	}
	return nil
}
