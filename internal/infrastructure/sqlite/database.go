package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const DefaultDatabaseFolder = "/database"
const DatabaseName = "main.db"
const MigrationPath = "migrations"

type Database struct {
	MainDirPath string
	fullDbPath  string
	Connection  *sql.DB
}

func (db *Database) Init() {
	db.prepareFullDbPath()
	db.checkAndCreateDbFolder()
	err := db.RunMigrations(MigrationPath)
	if err != nil {
		fmt.Println(err)
	}
}

func (db *Database) prepareFullDbPath() {
	db.fullDbPath = db.MainDirPath + DefaultDatabaseFolder
}

func (db *Database) checkAndCreateDbFolder() {
	if _, err := os.Stat(db.fullDbPath); os.IsNotExist(err) {
		os.Mkdir(db.fullDbPath, 0777)
	}

	dbFullPath := db.fullDbPath + "/" + DatabaseName
	if _, err := os.Stat(dbFullPath); os.IsNotExist(err) {
		os.Create(dbFullPath)
	}
}

func (db *Database) RunMigrations(migrationPath string) error {

	fmt.Println(db.fullDbPath + "/" + DatabaseName)
	connection, err := sql.Open("sqlite3", db.fullDbPath+"/"+DatabaseName)
	if err != nil {
		return fmt.Errorf("error while trying to connect to database: %v", err)
	}

	db.Connection = connection

	driver, err := sqlite3.WithInstance(connection, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error while trying to create database driver: %v", err)
	}
	fmt.Println(err)
	fmt.Println("file://" + db.MainDirPath + "/" + migrationPath)
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+db.MainDirPath+"/"+migrationPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error while trying to intitalize migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while executing migrations: %v", err)
	}

	log.Println("Successfully migrate !")
	return nil
}
