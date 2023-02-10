package dbsaver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type DBSaver struct {
	DB  *sql.DB
	ctx context.Context
	cfg *config.Config

	Buffer []encoder.Encode
}

func NewSaver(ctx context.Context, cfg *config.Config) (*DBSaver, error) {
	if cfg.DBAddress == "" {
		err := fmt.Errorf("there is no DB address")
		return nil, err
	}
	db, err := sql.Open("pgx", cfg.DBAddress)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &DBSaver{
		DB:  db,
		ctx: ctx,
		cfg: cfg,

		Buffer: make([]encoder.Encode, 0, 35),
	}, nil
}

func (s *DBSaver) CreateTable() {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*3)
	defer cancel()

	query := "CREATE TABLE IF NOT EXISTS metrics " +
		"(NAME VARCHAR(30), " +
		"TYPE VARCHAR(10), " +
		"VALUE DOUBLE PRECISION, " +
		"DELTA BIGINT, " +
		"HASH VARCHAR(100)" +
		");"
	_, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error during creating a new DB %s", err)
	}
}

func (s *DBSaver) SaveMetric(metric *encoder.Encode) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*3)
	defer cancel()

	query := "INSERT INTO metrics " +
		"(NAME, TYPE, VALUE, DELTA, HASH) " +
		"VALUES ($1, $2, $3, $4, $5)"

	_, err := s.DB.ExecContext(ctx, query,
		metric.ID, metric.MType, metric.Value, metric.Delta, metric.Hash)
	if err != nil {
		log.Printf("Error during insert a new DB %s", err)
		return err
	}
	log.Printf("metric %s send to the storage", metric.ID)
	return nil
}

func (s *DBSaver) SaveAllMetrics(metric encoder.Encode) error {
	s.Buffer = append(s.Buffer, metric)

	if cap(s.Buffer) == len(s.Buffer) {
		err := s.Flush()
		if err != nil {
			return errors.New("cannot add records to the database")
		}
	}
	return nil
}

func (s *DBSaver) Flush() error {
	if s.DB == nil {
		err := fmt.Errorf("you haven`t opened the database connection")
		return err
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO metrics " +
		"(NAME, TYPE, VALUE, DELTA, HASH) " +
		"VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Printf("statement error: %s", err)
		}
	}()

	for _, v := range s.Buffer {
		if _, err = stmt.Exec(v.ID, v.MType, float64(*v.Value), int(*v.Delta), v.Hash); err != nil {
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

	s.Buffer = s.Buffer[:0]
	return nil
}

func (s *DBSaver) LoadMetrics() ([]encoder.Encode, error) {
	var metrics []encoder.Encode

	if s.DB == nil {
		err := fmt.Errorf("you haven`t opened the database connection")
		return nil, err
	}

	rows, err := s.DB.QueryContext(s.ctx, "SELECT * FROM metrics")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var enc encoder.Encode
		err = rows.Scan(&enc.ID, &enc.MType, &enc.Value, &enc.Delta, &enc.Hash)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, enc)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (s *DBSaver) KeeperCheck() error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*3)
	defer cancel()
	err := s.DB.PingContext(ctx)
	return err
}
