package zt_formatter

import (
	"fmt"
	"os"
	"path"
	"runtime"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
)


func ExampleFormatter_Format_IncludingFooAndFile() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)

	l.SetReportCaller(true)
	l.SetFormatter(&zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: formatter.Formatter{
			HideKeys: true,
			FieldsOrder: []string{"component", "category"},
		},
	})

	l.Debug("test1")
	l.Info("test2")
	l.Warn("test3")
	l.Error("test4")

	// Output:
	// - [DEBU] test1
	// - [INFO] test2
	// - [WARN] test3
	// - [ERRO] test4
}


func ExampleFormatter_Format_default() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
	},
	})

	l.Debug("test1")
	l.Info("test2")
	l.Warn("test3")
	l.Error("test4")

	// Output:
	// - [DEBU] test1
	// - [INFO] test2
	// - [WARN] test3
	// - [ERRO] test4
}

func ExampleFormatter_Format_full_level() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		ShowFullLevel:   true,
	},
	})

	l.Debug("test1")
	l.Info("test2")
	l.Warn("test3")
	l.Error("   test4")

	// Output:
	// - [DEBUG] test1
	// - [INFO] test2
	// - [WARNING] test3
	// - [ERROR]    test4
}
func ExampleFormatter_Format_show_keys() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        false,
	},
	})

	ll := l.WithField("category", "rest")

	l.Info("test1")
	ll.Info("test2")

	// Output:
	// - [INFO] test1
	// - [INFO] [category:rest] test2
}

func ExampleFormatter_Format_hide_keys() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        true,
	},
	})

	ll := l.WithField("category", "rest")

	l.Info("test1")
	ll.Info("test2")

	// Output:
	// - [INFO] test1
	// - [INFO] [rest] test2
}

func ExampleFormatter_Format_sort_order() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		HideKeys:        false,
	},
	})

	ll := l.WithField("component", "main")
	lll := ll.WithField("category", "rest")

	l.Info("test1")
	ll.Info("test2")
	lll.Info("test3")

	// Output:
	// - [INFO] test1
	// - [INFO] [component:main] test2
	// - [INFO] [category:rest] [component:main] test3
}

func ExampleFormatter_Format_field_order() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		NoColors:        true,
		TimestampFormat: "-",
		FieldsOrder:     []string{"component", "category"},
		HideKeys:        false,
	},
	})

	ll := l.WithField("component", "main")
	lll := ll.WithField("category", "rest")

	l.Info("test1")
	ll.Info("test2")
	lll.Info("test3")

	// Output:
	// - [INFO] test1
	// - [INFO] [component:main] test2
	// - [INFO] [component:main] [category:rest] test3
}

func ExampleFormatter_Format_trim_message() {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&zt_formatter.ZtFormatter{Formatter:formatter.Formatter{
		TrimMessages:    true,
		NoColors:        true,
		TimestampFormat: "-",
	},
	})

	l.Debug(" test1 ")
	l.Info("test2 ")
	l.Warn(" test3")
	l.Error("   test4   ")

	// Output:
	// - [DEBU] test1
	// - [INFO] test2
	// - [WARN] test3
	// - [ERRO] test4
}
