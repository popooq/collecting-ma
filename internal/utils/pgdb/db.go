package pgdb

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/storage"
)

type DataBase struct {
	DB  *sql.DB
	ctx context.Context
	cfg *config.Config
	str *storage.MetricsStorage
}

func New(ctx context.Context, cfg *config.Config, str *storage.MetricsStorage) *DataBase {
	if cfg.DBAddress == "" {
		return nil
	}
	db, err := sql.Open("pgx", cfg.DBAddress)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}
	return &DataBase{
		DB:  db,
		ctx: ctx,
		cfg: cfg,
		str: str,
	}
}

func (db *DataBase) CreateTable() {
	ctx, cancel := context.WithTimeout(db.ctx, time.Second*3)
	defer cancel()
	db.DB.ExecContext(ctx, "CREATE TABLE metrics (ID SERIAL PRIMARY KEY, "+
		"NAME CHARACTER VARYING(30), "+
		"TYPE CHARACTER VARYING(10), "+
		"HASH CHARACTER VARYING(100), "+
		"VALUE DOUBLE PRECISION, "+
		"DELTA INTEGER"+
		");")
}

func (db *DataBase) ReturnCntext() context.Context {
	return db.ctx
}
