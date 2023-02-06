package pgdb

import (
	"context"
	"database/sql"
	"fmt"
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

	Buffer []encoder.Encode
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

		Buffer: make([]encoder.Encode, 0, 35),
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
	log.Printf("metric %s send to the storage", enc.ID)
}

func (db *DataBase) Flush() error {
	// проверим на всякий случай
	if db.DB == nil {
		err := fmt.Errorf("you haven`t opened the database connection")
		return err
	}
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO metrics " +
		"(NAME, TYPE, HASH, VALUE, DELTA) " +
		"VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range db.Buffer {
		if _, err = stmt.Exec(v.ID, v.MType, v.Hash, v.Value, v.Delta); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("update drivers: unable to rollback: %v", err)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("update drivers: unable to commit: %v", err)
		return err
	}

	db.Buffer = db.Buffer[:0]
	return nil
}

func (db *DataBase) ReturnCntext() context.Context {
	return db.ctx
}
