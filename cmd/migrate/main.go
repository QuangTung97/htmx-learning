package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"htmx/config"
	"htmx/config/prod"
	"htmx/pkg/migration"
)

func main() {
	unv := prod.NewUniverse()
	conf := config.Loc.Get(unv)

	cmd := migration.MigrateCommand(conf.MySQL.DSN())
	err := cmd.Execute()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
