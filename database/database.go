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
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"user_name"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

func (c *Cfg) Default() {

}

func (c *Cfg) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", c.Host, c.Username, c.Password, c.DBName, c.Port)
}

type Database interface {
	DB() *gorm.DB
}

type DatabaseImpl struct {
	db *gorm.DB
	l  *mxlog.Logger
}

type Params struct {
	fx.In

	Lc     fx.Lifecycle
	Log    mxlog.Loggers `name:"mxlog"`
	Config config.Config `name:"config"`
}

type Result struct {
	fx.Out

	Database Database
}

func NewDatabase(p Params) Result {
	c := p.Config.GetItem("database").(Cfg)

	impl := &DatabaseImpl{
		l: p.Log.GetLogger("databse"),
	}

	p.Lc.Append(fx.Hook{
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

	return Result{Database: impl}
}

func (impl DatabaseImpl) DB() *gorm.DB {
	return impl.db
}

var Module = fx.Options(fx.Provide(NewDatabase))
