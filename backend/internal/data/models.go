package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users UserModel
}

func NewModels(dbpool *pgxpool.Pool) Models {
	return Models{
		Users: UserModel{dbpool: dbpool},
	}
}
