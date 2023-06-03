package database

import (
	"context"

	"github.com/dafsic/go-pkgs/config"
	"github.com/dafsic/go-pkgs/mxlog"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Params struct {
	fx.In

	Lc     fx.Lifecycle
	Log    mxlog.Loggers
	Config config.Config
}

type Result struct {
	fx.Out

	Database Database `name:"postgres"`
}

func NewPostgresDatabase(p Params) Result {
	c := p.Config.GetItem(ModuleName).(Cfg)

	impl := &DatabaseImpl{
		l: p.Log.GetLogger(ModuleName),
	}

	p.Lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			var err error
			impl.db, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  c.DSN,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			return err
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			// gorm维护连接池，不用关闭
			return nil
		},
	})

	return Result{Database: impl}
}

func (impl DatabaseImpl) Inst() *gorm.DB {
	return impl.db
}

var Module = fx.Options(fx.Provide(NewPostgresDatabase))
