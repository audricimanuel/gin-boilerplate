package tools

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewGormDB(ctx context.Context, dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error when NewGormDB, error: %s\n", err.Error())
		return nil
	}

	return db.WithContext(ctx)
}
