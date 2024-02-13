package migration

// MUST import in main.go to run
// _ "github.com/golang-migrate/migrate/v4/database/mysql"
// _ "github.com/golang-migrate/migrate/v4/source/file"
import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

const versionTimeFormat = "20060102150405"

func migrateUpCommand(sourceURL, databaseURL string) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "migrate all the way up",
		Run: func(cmd *cobra.Command, args []string) {
			m, err := migrate.New(sourceURL, databaseURL)
			if err != nil {
				panic(err)
			}

			err = m.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No change in migration")
				return
			}
			if err != nil {
				panic(err)
			}

			fmt.Println("Migrated up")
		},
	}
}

func migrateDownCommand(sourceURL, databaseURL string) *cobra.Command {
	return &cobra.Command{
		Use:   "down [number]",
		Short: "migrate down by 'number' times",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}

			m, err := migrate.New(sourceURL, databaseURL)
			if err != nil {
				panic(err)
			}

			err = m.Steps(-num)
			if err != nil {
				panic(err)
			}

			fmt.Println("Migrated down ", num)
		},
	}
}

func migrateForceCommand(sourceURL, databaseURL string) *cobra.Command {
	return &cobra.Command{
		Use:   "force [version]",
		Short: "force dirty migration using 'version'",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			m, err := migrate.New(sourceURL, databaseURL)
			if err != nil {
				panic(err)
			}

			version, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}

			err = m.Force(version)
			if err != nil {
				panic(err)
			}

			fmt.Println("Forced version:", version)
		},
	}
}

func migrateCreateCommand(migrationDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "create [name]",
		Short: "create a SQL migration script with format {timestamp}_{name}",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			version := now.Format(versionTimeFormat)
			name := args[0]

			up := fmt.Sprintf("%s/%s_%s.up.sql", migrationDir, version, name)
			down := fmt.Sprintf("%s/%s_%s.down.sql", migrationDir, version, name)

			err := os.WriteFile(up, []byte{}, 0644)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(down, []byte{}, 0644)
			if err != nil {
				panic(err)
			}

			fmt.Println("Created SQL up script:", up)
			fmt.Println("Created SQL down script:", down)
		},
	}
}

const migrationDirectory = "migrations"

// MigrateCommand the command for migration
func MigrateCommand(dsn string) *cobra.Command {
	databaseURL := fmt.Sprintf("mysql://%s", dsn)
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "database migration command",
	}

	sourceURL := fmt.Sprintf("file://%s", migrationDirectory)

	fmt.Println("Source URL:", sourceURL)
	fmt.Println("------------------------------------------------------------")

	cmd.AddCommand(
		migrateUpCommand(sourceURL, databaseURL),
		migrateDownCommand(sourceURL, databaseURL),
		migrateForceCommand(sourceURL, databaseURL),
		migrateCreateCommand(migrationDirectory),
	)

	return cmd
}

// MigrateUpForTesting ...
func MigrateUpForTesting(rootDir string, dsn string) {
	sourceURL := fmt.Sprintf("file://%s", path.Join(rootDir, migrationDirectory))
	databaseURL := fmt.Sprintf("mysql://%s", dsn)

	fmt.Println("SourceURL:", sourceURL)
	// fmt.Println("DatabaseURL:", databaseURL)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No change in migration")
		return
	}
	if err != nil {
		panic(err)
	}
}
