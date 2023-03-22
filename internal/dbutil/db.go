package dbutil

import (
	"database/sql"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Transactor interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}

func IsUniqueViolation(err error, constraint string) bool {
	// NOTE: How to check if error type is DUPLICATE KEY in GORM.
	// https://github.com/go-gorm/gorm/issues/4037
	pgError, ok := err.(*pgconn.PgError)
	return ok && pgError.Code == "23505" && pgError.ConstraintName == constraint
}
