// robot 参数读写,复杂后可以引入gorm
package database

import (
	"github.com/dafsic/go-pkgs/mxlog"
	"gorm.io/gorm"
)

const ModuleName = "database"

type Cfg struct {
	DSN string `toml:"dsn"`
}

func (c *Cfg) Default() {
	c.DSN = ""
}

type Database interface {
	Inst() *gorm.DB
}

type DatabaseImpl struct {
	db *gorm.DB
	l  *mxlog.Logger
}
