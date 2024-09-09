package database

import (
	"context"
	"fmt"
	"gin-boilerplate/src/config"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type (
	DBCollection struct {
		MongoDB        *mongo.Database
		PostgresDBSqlx *sqlx.DB
		PostgresDBGorm *gorm.DB
	}
)

func NewDatabaseCollection(cfg config.Config) DBCollection {
	ctx := context.Background()

	// mongodb
	mongoDBConfig := cfg.DataSource.MongoDBConfig
	mongoDB := InitializeMongoDatabase(ctx, mongoDBConfig.ConnectionString, mongoDBConfig.DatabaseName)

	// postgres
	postgresDBConfig := cfg.DataSource.PostgresDBConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		postgresDBConfig.Host, postgresDBConfig.User, postgresDBConfig.Password, postgresDBConfig.Name, postgresDBConfig.Port, postgresDBConfig.SSLMode, postgresDBConfig.Timezone,
	)

	// postgres with sqlx
	postgresDBSqlx := InitializePostgresqlDatabaseSqlx(ctx, dsn)
	// postgres with gorm
	postgresDBGorm := InitializePostgresqlDatabaseGorm(ctx, dsn)

	return DBCollection{
		MongoDB:        mongoDB,
		PostgresDBSqlx: postgresDBSqlx,
		PostgresDBGorm: postgresDBGorm,
	}
}
