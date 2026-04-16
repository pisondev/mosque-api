package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func ConnectDB(log *logrus.Logger) *pgxpool.Pool {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("database url is not set in env")
	}

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("unable to parse database config: %v", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	log.Info("database connected successfully")
	return db
}
