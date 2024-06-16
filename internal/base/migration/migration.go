package migration

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationConfig struct {
	//db is the database instance
	Db *sql.DB

	//Dialect is the type of Database (mysql)
	Dialect string

	//MigrationDir is the location of migration folder (default to db/migrations)
	MigrationDir string
}

func NewMigrationConfig(migrationDir, dbHost, dbPort, dbUser, dbPass, dbName, dbDriver string) (*MigrationConfig, error) {
	log.Info().Msg("db Host : " + dbHost)
	log.Info().Msg("db Port : " + dbPort)
	log.Info().Msg("db name : " + dbName)
	log.Info().Msg("db driver : " + dbDriver)

	migrationConf := MigrationConfig{MigrationDir: migrationDir}
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", dbUser, dbPass, dbHost, dbPort, dbName)
	//dbDriver := os.Getenv("DB_DRIVER")

	var err error
	//open db connection based on driver
	switch dbDriver {
	case "mysql":
		migrationConf.Dialect = dbDriver
		migrationConf.Db, err = sql.Open(dbDriver, dbDSN)
		if err != nil {
			return nil, errors.New("Migrate.NewMigrationConfig err : " + err.Error())
		}
	default:
		return nil, errors.New("error db driver is not found (currently mysql supported only)")
	}

	return &migrationConf, nil
}

//MigrateUp - will migrate the database to the latest version
func MigrateUp(config *MigrationConfig) error {
	log.Info().Msg("Migrating up database ...")
	driver, errDriver := mysql.WithInstance(config.Db, &mysql.Config{})
	if errDriver != nil {
		return errors.New("Migrate.MigrateUp errDriver : " + errDriver.Error())
	}

	migrateDatabase, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationDir),
		config.Dialect, driver)
	if err != nil {
		return errors.New("Migrate.MigrateUp err : " + err.Error())
	}

	errUp := migrateDatabase.Up()
	if errUp != nil {
		return errUp
	}

	log.Info().Msg("Migration done ...")

	//get latest version
	version, dirty, errVersion := migrateDatabase.Version()
	//ignore error in this line. Skip the version check
	if errVersion != nil {
		return errors.New("Migrate.MigrateUp errVersion : " + errVersion.Error())
	}

	if dirty {
		log.Info().Msg("dirty migration. Please clean up database")
	}

	msgLatestVersion := fmt.Sprintf("latest version is %d", version)
	log.Info().Msg(msgLatestVersion)
	return nil
}

//MigrateDown - will migrate the database to the latest version
func MigrateDown(config *MigrationConfig) error {
	log.Info().Msg("Migrating down database ...")
	driver, errDriver := mysql.WithInstance(config.Db, &mysql.Config{})
	if errDriver != nil {
		return errors.New("Migrate.MigrateDown errDriver : " + errDriver.Error())
	}

	migrateDatabase, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationDir),
		config.Dialect, driver)
	if err != nil {
		return errors.New("Migrate.MigrateDown err : " + err.Error())
	}

	errDown := migrateDatabase.Down()
	if errDown != nil {
		return errDown
	}

	log.Info().Msg("Migration done ...")

	return nil
}
