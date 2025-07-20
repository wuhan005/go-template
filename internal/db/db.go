// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/wuhan005/go-template/internal/conf"
	"github.com/wuhan005/go-template/internal/dbutil"
)

var tables = []interface{}{
	&User{},
}

var dbInstance *gorm.DB

// Init initializes the database.
func Init() (*gorm.DB, error) {
	dsn := conf.Postgres.DSN
	dsnURL, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "parse dsn")
	}

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

	if err := db.Use(tracing.NewPlugin(
		tracing.WithAttributes(
			attribute.String("db.name", db.Name()),
			attribute.String("db.ip", fmt.Sprintf("%s:%d", dsnURL.Host, dsnURL.Port)),
		),
	)); err != nil {
		return nil, errors.Wrap(err, "register otelgorm plugin")
	}

	// Migrate databases.
	if err := db.AutoMigrate(tables...); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	SetDatabaseStore(db)

	dbInstance = db

	return db, nil
}

// Ping checks the database connection.
func Ping(ctx context.Context) error {
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("get db: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// SetDatabaseStore sets the database table store.
func SetDatabaseStore(db *gorm.DB) {
	Users = NewUsersStore(db)
}
