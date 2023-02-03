package backuper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/pgdb"
)

type Backuper struct {
	storage *storage.MetricsStorage
	cfg     *config.Config
	enc     *encoder.Encode
	DB      *pgdb.DataBase

	file   *os.File
	writer *bufio.Writer
}

func NewSaver(storage *storage.MetricsStorage, cfg *config.Config, enc *encoder.Encode, DB *pgdb.DataBase) (*Backuper, error) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return nil, err
	}

	return &Backuper{
		storage: storage,
		cfg:     cfg,
		enc:     enc,
		DB:      DB,

		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

type Loader struct {
	storage *storage.MetricsStorage
	cfg     *config.Config
	encoder *encoder.Encode

	file   *os.File
	reader *bufio.Reader
}

func NewLoader(storage *storage.MetricsStorage, cfg *config.Config, encoder *encoder.Encode) (*Loader, error) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}

	return &Loader{
		storage: storage,
		cfg:     cfg,
		encoder: encoder,

		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (s *Backuper) Close() error {
	return s.file.Close()
}

func (s *Backuper) SaveToFile() error {
	err := s.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("err %w", err)
	}

	_ = s.writer.WriteByte('[')

	for k, v := range s.storage.MetricsGauge {
		v := v

		s.fillEncGauge(k, "gauge", v)

		data, err := s.enc.Marshall()
		if err != nil {
			return err
		}

		_, err = s.writer.Write(data)
		if err != nil {
			return err
		}

		err = s.writer.WriteByte(',')
		if err != nil {
			return err
		}

		err = s.writer.WriteByte('\n')
		if err != nil {
			return err
		}
	}

	for k, d := range s.storage.MetricsCounter {
		d := d

		s.fillEncCounter(k, "counter", d)

		data, err := s.enc.Marshall()
		if err != nil {
			return err
		}

		_, err = s.writer.Write(data)
		if err != nil {
			return err
		}

		err = s.writer.WriteByte(',')
		if err != nil {
			return err
		}

		err = s.writer.WriteByte('\n')
		if err != nil {
			return err
		}
	}

	_, _ = s.writer.WriteString("{}]")
	log.Print("new save to file")
	return s.writer.Flush()
}

func (s *Backuper) SaveToDB() error {
	for k, v := range s.storage.MetricsGauge {
		v := v
		log.Println(k, v)
		s.fillEncGauge(k, "gauge", v)

		s.DB.InsertMetric(*s.enc)
	}

	for k, d := range s.storage.MetricsCounter {
		d := d
		log.Println(k, d)

		s.fillEncCounter(k, "counter", d)

		s.DB.InsertMetric(*s.enc)
	}

	log.Print("new save to DB")

	err := fmt.Errorf("error")
	return err

}

func (s *Backuper) Saver() {
	tickerstore := time.NewTicker(s.cfg.StoreInterval)

	for {
		<-tickerstore.C

		if s.cfg.DBAddress != "" {

			s.DB.CreateTable()

			err := s.SaveToDB()
			if err != nil {
				return
			}

		}
		err := s.SaveToFile()
		if err != nil {
			return
		}
	}
}

func (l *Loader) Close() error {
	return l.file.Close()
}

func (l *Loader) LoadFromFile() error {
	var (
		data    []byte
		encoder []encoder.Encode
	)

	data, err := io.ReadAll(l.reader)
	if err != nil {
		log.Printf("erad err : %s", err)
		return err
	}

	err = json.Unmarshal(data, &encoder)
	if err != nil {
		log.Printf("marshal err : %s", err)
		return err
	}

	for _, v := range encoder {
		switch v.MType {
		case "gauge":
			l.storage.GetBackupGauge(v.ID, *v.Value)
		case "counter":
			l.storage.GetBackupCounter(v.ID, *v.Delta)
		}
	}

	return nil
}

func (s *Backuper) fillEncGauge(k, mtype string, v float64) {
	s.enc.ID = k
	s.enc.MType = mtype
	s.enc.Value = &v
}

func (s *Backuper) fillEncCounter(k, mtype string, d int64) {
	s.enc.ID = k
	s.enc.MType = mtype
	s.enc.Delta = &d
}
