package db

import (
	"context"
)

type Repository interface {
	GetOne(context.Context, string, interface{}, ...interface{}) error
	GetAll(context.Context, string, interface{}, ...interface{}) error
}
