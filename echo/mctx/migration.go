package mctx

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
)

// https://github.com/golang-migrate/migrate/blob/master/testing/testing.go
type Migration struct {
	Migrate *migrate.Migrate
	Driver  string
}

// Up is func new version
func (mi *Migration) Up() (bool, error) {
	err := mi.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("ErrNoChange: ", err.Error())
			return true, nil
		}
		return false, err
	}
	return true, nil
}

// Down is func tear down version
func (mi *Migration) Down() (bool, error) {
	err := mi.Migrate.Down()
	if err != nil {
		return false, err
	}
	return true, nil
}

// RunMigration is func run migration mysql
func RunMigration(dbConn *sql.DB, migrationsFolderLocation string) (*Migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	_, err := url.Parse(pathToMigrate)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	driver, err := _mysql.WithInstance(dbConn, &_mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "mysql", driver)
	if err != nil {
		return nil, err
	}
	return &Migration{Migrate: m}, nil
}
