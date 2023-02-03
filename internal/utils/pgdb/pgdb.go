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
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
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

	query := "CREATE TABLE metrics " +
		"(NAME VARCHAR(30), " +
		"TYPE VARCHAR(10), " +
		"HASH VARCHAR(100), " +
		"VALUE DOUBLE PRECISION, " +
		"DELTA BIGINT" +
		");"
	_, err := db.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error during creating a new DB %s", err)
	}
}

func (db *DataBase) TruncateMetric() {
	ctx, cancel := context.WithTimeout(db.ctx, time.Second*3)
	defer cancel()

	query := "TRUNCATE TABLE metrics;"

	_, err := db.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error during truncate table %s", err)
	}
}

func (db *DataBase) InsertMetric(enc encoder.Encode) {
	ctx, cancel := context.WithTimeout(db.ctx, time.Second*3)
	defer cancel()

	query := "INSERT INTO metrics " +
		"(NAME, TYPE, HASH, VALUE, DELTA) " +
		"VALUES ($1, $2, $3, $4, $5)"

	_, err := db.DB.ExecContext(ctx, query,
		enc.ID, enc.MType, enc.Hash, enc.Value, enc.Delta)
	if err != nil {
		log.Printf("Error during insert a new DB %s", err)
	}
}

func (db *DataBase) ReturnCntext() context.Context {
	return db.ctx
}
