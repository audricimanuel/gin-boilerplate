package tools

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

func NewSqlxDsn(ctx context.Context, driver, dsn string) *sqlx.DB {
	log := logrus.WithContext(ctx)

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Fatalf("error when NewSqlxDsn, error: %s\n", err.Error())
		return nil
	}

	// Setup Connection
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("ping %s", driver)
	if err := db.Ping(); err != nil {
		log.Fatalf("error when ping %s, error: %s\n", driver, err.Error())
	}

	return db
}
