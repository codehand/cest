package mctx

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_mysql "github.com/golang-migrate/migrate/v4/database/mysql"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// https://github.com/golang-migrate/migrate/blob/master/testing/testing.go
type Migration struct {
	Migrate *migrate.Migrate
	Driver  string
}

// Up is func new version
func (mi *Migration) Up() (bool, error) {
	log.Println("Starting a Test. Migrating the Database")
	err := mi.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Migrating ErrNoChange: ", err.Error())
			return true, nil
		}
		log.Printf("Migrating err %v\n", err)
		return false, err
	}
	log.Println("Database Migrated Successfully")
	return true, nil
}

// Down is func tear down version
func (mi *Migration) Down() (bool, error) {
	log.Println("Finishing Test. Dropping The Database")
	err := mi.Migrate.Down()
	if err != nil {
		log.Printf("Migrating err %v\n", err)
		return false, err
	}
	log.Println("Database Dropped Successfully")
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
