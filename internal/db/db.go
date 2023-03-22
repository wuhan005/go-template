package db

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/wuhan005/go-template/internal/dbutil"
)

var AllTables = []interface{}{
	// TODO ...
}

// Init initializes the database.
func Init() (*gorm.DB, error) {
	dsn := os.ExpandEnv("postgres://$PGUSER:$PGPASSWORD@$PGHOST/$PGNAME?sslmode=$PGSSLMODE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return dbutil.Now()
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             3 * time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "open connection")
	}

	// Migrate databases.
	if err := db.AutoMigrate(AllTables...); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	// Create sessions table.
	q := `
CREATE TABLE IF NOT EXISTS sessions (
    key        TEXT PRIMARY KEY,
    data       BYTEA NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL
);`
	if err := db.Exec(q).Error; err != nil {
		return nil, errors.Wrap(err, "create sessions table")
	}

	SetDatabaseStore(db)

	return db, nil
}

// SetDatabaseStore sets the database table store.
func SetDatabaseStore(db *gorm.DB) {
	// TODO ...
}
