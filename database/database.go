// robot 参数读写,复杂后可以引入gorm
package database

import (
	"context"
	"fmt"

	"github.com/dafsic/go-pkgs/config"
	"github.com/dafsic/go-pkgs/mxlog"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Cfg struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func (c Cfg) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", c.Host, c.Username, c.Password, c.DBName, c.Port)
}

type Database interface {
	DB() *gorm.DB
}

type DatabaseImpl struct {
	db *gorm.DB
	l  *mxlog.Logger
}

func NewDatabase(lc fx.Lifecycle, log mxlog.Loggers, cfg config.Config) Database {
	c := cfg.GetElem("database").(Cfg)

	impl := &DatabaseImpl{
		l: log.GetLogger("databse"),
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			var err error
			impl.db, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  c.DSN(),
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{})
			return err
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			// gorm维护连接池，不用关闭
			return nil
		},
	})

	return impl
}

func (impl DatabaseImpl) DB() *gorm.DB {
	return impl.db
}

var StoreModule = fx.Options(fx.Provide(NewDatabase))
