// robot 参数读写,复杂后可以引入gorm
package database

import (
	"fmt"

	"github.com/dafsic/go-pkgs/mxlog"
	"gorm.io/gorm"
)

const ModuleName = "database"

type Cfg struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"user_name"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

func (c *Cfg) Default() {}

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
