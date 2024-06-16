package cmd

import (
	"errors"
	"flag"
	"os"

	"store-api/internal/base/migration"

	"github.com/joho/godotenv"
	gologger "github.com/mo-taufiq/go-logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

const (
    MigrateUp = iota
    MigrateDown
)

var migrateUpCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate Up DB Paylater",
    Long:  `Please you know what are you doing by using this command`,
    Run: func(cmd *cobra.Command, args []string) {
        envName, _ := cmd.Flags().GetString("env")
        loadEnv(envName)

        runMigration(MigrateUp)
    },
}

var migrateDownCmd = &cobra.Command{
    Use:   "migratedown",
    Short: "Migrate Down DB Paylater",
    Long:  `Please you know what are you doing by using this command`,
    Run: func(cmd *cobra.Command, args []string) {
        envName, _ := cmd.Flags().GetString("env")
        loadEnv(envName)

        runMigration(MigrateDown)
    },
}

func init() {
    rootCmd.AddCommand(migrateUpCmd)
    rootCmd.AddCommand(migrateDownCmd)

    migrateUpCmd.PersistentFlags().StringP("env", "e", "prod", "environment type")
    migrateDownCmd.PersistentFlags().StringP("env", "e", "prod", "environment type")
}

func runMigration(direction int) {
    pathMigration := os.Getenv("APP_MIGRATION_PATH")
    migrationDir := flag.String("migration-dir", pathMigration, "migration directory")
    log.Info().Msg("path migration : " + pathMigration)

    migrationConf, errMigrationConf := migration.NewMigrationConfig(*migrationDir,
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USERNAME"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        "mysql")
    if errMigrationConf != nil {
        log.Error().Msg(errMigrationConf.Error())
        return
    }

    var errMigration error
    switch direction {
    case MigrateUp:
        errMigration = migration.MigrateUp(migrationConf)
        break
    case MigrateDown:
        errMigration = migration.MigrateDown(migrationConf)
        break
    default:
        errMigration = errors.New("Unknown migration direction")
    }
    if errMigration != nil {
        if errMigration.Error() != "no change" {
            log.Error().Msg(errMigration.Error())
            return
        }
        log.Info().Msg("Migration success : no change table . . .")
    }
}

func loadEnv(envName string) {
    gologger.LogConf.NestedLocationLevel = 2
    log.Logger = log.Output(
        zerolog.ConsoleWriter{
            Out:     os.Stderr,
            NoColor: false,
        },
    )

    dotenvPath := "params/.env"

    if envName == "test" {
        dotenvPath = "params/.env.test"
    }

    err := godotenv.Load(dotenvPath)
    if err != nil {
        log.Error().Msg("Error loading .env file")
    }
}
