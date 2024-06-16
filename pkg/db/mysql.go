package db

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"

	// mysql driver
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MySQLClientRepository struct {
	DB *sqlx.DB
	TZ string
}

// NewMySQLClient creates DB client
func NewMySQLRepository(host, uname, pass, dbname string, port int) (*MySQLClientRepository, error) {
	tz := "&loc=Asia%2FJakarta"

	if os.Getenv("DB_TZ") == "UTC" { // First phase, no set DB_TZ to get Asia Jakarta
		tz = "" // Second phase, split db, set UTC for MariaDB
	}

	sqltrace.Register("mysql", &mysql.MySQLDriver{})

	dsnFormat := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true%s", uname, pass, host, port, dbname, tz)
	db, err := sqlxtrace.Connect("mysql", dsnFormat)
	if err != nil {
		logrus.Error(fmt.Sprintf("Cannot connect to MySQL. %v", err))
		return nil, errors.Wrap(err, "Cannot connect to MySQL")
	}
	if db == nil {
		panic("missing db")
	}

	return &MySQLClientRepository{DB: db, TZ: tz}, nil
}

//Ref https://github.com/jmoiron/sqlx
func (r MySQLClientRepository) GetOne(ctx context.Context, query string, target interface{}, args ...interface{}) error {
	err := r.DB.Get(target, query, args...)
	if err != nil {
		return err
	}
	return err
}

func (r MySQLClientRepository) GetAll(ctx context.Context, query string, target interface{}, args ...interface{}) error {
	err := r.DB.Select(target, query, args...)
	if err != nil {
		return err
	}
	return err
}
