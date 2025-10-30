package gorm

import (
	"astroneko-backend/internal/core/domain/shared"
	"fmt"

	// Import postgres driver for side effects
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB struct
type DB struct {
	Postgres *gorm.DB
}

var dbConnect = &DB{}

// ConnectToPostgreSQL connects to PostgreSQL database
func ConnectToPostgreSQL(host, port, username, pass, dbname string, sslmode bool) (*DB, error) {
	// config string
	var connectionStr string

	if host == "" && port == "" && dbname == "" {
		return nil, shared.ErrDatabaseConnectionFailed
	}

	if port == "APP_DATABASE_POSTGRES_PORT" {
		port = "5432"
	}

	if sslmode {
		connectionStr = fmt.Sprintf(
			"host=%v user=%v password=%v dbname=%v port=%v sslmode=require TimeZone=UTC",
			host, username, pass, dbname, port)
	} else {
		connectionStr = fmt.Sprintf(
			"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=UTC",
			host, username, pass, dbname, port)
	}

	// connect postgres
	dial := postgres.Open(connectionStr)

	var err error
	pg, err := gorm.Open(dial, &gorm.Config{
		DryRun: false,
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	dbConnect.Postgres = pg

	return dbConnect, nil
}

func ConnectToCloudSQL(instanceConnectionName, dbUser, dbPassword, dbName string) (*DB, error) {
	connectionStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable TimeZone=UTC",
		dbUser, dbPassword, dbName, instanceConnectionName,
	)

	dial := postgres.Open(connectionStr)

	pg, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	dbConnect.Postgres = pg
	return dbConnect, nil
}

// DisconnectPostgres disconnects from PostgreSQL database
func DisconnectPostgres(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic("close db")
	}
	err = sqlDb.Close()
	if err != nil {
		panic("close db")
	}
	logrus.Println("Connected with postgres has closed")
}
