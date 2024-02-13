package config

import (
	"fmt"
	"net/url"

	"github.com/QuangTung97/svloc"
	"github.com/jmoiron/sqlx"
)

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         uint16 `mapstructure:"port"`
	Database     string `mapstructure:"database"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Options      string `mapstructure:"options"`
}

// DSN returns data source name
func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		url.PathEscape(c.Username),
		url.PathEscape(c.Password),
		c.Host, c.Port, c.Database, c.Options,
	)
}

// MustConnect connects to database using sqlx
func (c MySQLConfig) MustConnect() *sqlx.DB {
	db := sqlx.MustConnect("mysql", c.DSN())

	fmt.Println("MaxOpenConns:", c.MaxOpenConns)
	fmt.Println("MaxIdleConns:", c.MaxIdleConns)
	fmt.Println("Options:", c.Options)

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	return db
}

var DBLoc = svloc.Register[*sqlx.DB](func(unv *svloc.Universe) *sqlx.DB {
	return Loc.Get(unv).MySQL.MustConnect()
})
