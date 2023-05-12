// 日志的level等信息应该从配置文件里读取,转储可以用linux的logrotate
package mxlog

import (
	"io"
	"os"
	"sync"

	"github.com/dafsic/go-pkgs/config"
	"go.uber.org/fx"
)

type Cfg struct {
	Level string
}

type Loggers interface {
	GetLogger(name string) *Logger
}

type LoggersImpl struct {
	mux     sync.Mutex
	lvl     string
	output  io.Writer
	loggers map[string]*Logger
}

func (l *LoggersImpl) GetLogger(name string) *Logger {
	l.mux.Lock()
	defer l.mux.Unlock()
	i, ok := l.loggers[name]
	if !ok {
		i = NewLogger(l.output, name, LogLevelFromString(l.lvl), Ldefault)
		l.loggers[name] = i
	}
	return i
}

func NewLoggers(c config.Config) Loggers {
	cfg := c.GetElem("mxlog").(Cfg)
	t := &LoggersImpl{
		output:  os.Stdout,
		lvl:     cfg.Level,
		loggers: make(map[string]*Logger, 8),
	}

	return t
}

var Module = fx.Options(fx.Provide(NewLoggers))

// var once sync.Once
// var l Loggers

// func NewLoggerslog(cfg config.Config) Loggers {
// 	once.Do(func() {
// 		var t LoggersImpl
// 		t.output = os.Stdout
// 		t.lvl = cfg.GetElem("logLevel").(string)
// 		t.loggers = make(map[string]*Logger, 8)
// 		l = &t
// 	})
// 	return l
// }
