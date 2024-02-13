package integration

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/QuangTung97/svloc"
	"github.com/jmoiron/sqlx"

	"htmx/config"
	"htmx/model"
	"htmx/pkg/dbtx"
	"htmx/pkg/migration"

	// for integration test, must not be imported in any main.go
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestCase struct {
	Unv *svloc.Universe
}

var initOnce sync.Once

var globalConf config.Config
var globalDB *sqlx.DB

func NewTestCase() *TestCase {
	initOnce.Do(func() {
		rootDir := findRootDir()

		conf := config.LoadTestConfig(rootDir)
		migration.MigrateUpForTesting(rootDir, conf.MySQL.DSN())

		db := conf.MySQL.MustConnect()

		globalConf = conf
		globalDB = db
	})

	tc := &TestCase{
		Unv: svloc.NewUniverse(),
	}

	config.Loc.MustOverride(tc.Unv, globalConf)
	config.DBLoc.MustOverride(tc.Unv, globalDB)

	return tc
}

func (tc *TestCase) TruncateTables(tables ...model.GetTableName) {
	db := config.DBLoc.Get(tc.Unv)
	for _, table := range tables {
		fmt.Println("TRUNCATING Table:", table.TableName())
		db.MustExec(fmt.Sprintf("TRUNCATE %s", table.TableName()))
	}
}

func (tc *TestCase) Autocommit() context.Context {
	return dbtx.ProviderLoc.Get(tc.Unv).Autocommit(context.Background())
}

func findRootDir() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	directory := workdir
	for {
		files, err := os.ReadDir(directory)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == "go.mod" {
				return directory
			}
		}

		directory = path.Dir(directory)
	}
}
