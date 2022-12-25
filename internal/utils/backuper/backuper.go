package backuper

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
	"github.com/popooq/collectimg-ma/internal/utils/env"
	"github.com/popooq/collectimg-ma/internal/utils/storage"
)

type Backuper struct {
	storage *storage.MetricsStorage
	env     *env.ConfigServer
	enc     *encoder.Metrics

	file   *os.File
	writer *bufio.Writer
}

func NewSaver(storage *storage.MetricsStorage, env *env.ConfigServer, enc *encoder.Metrics) (*Backuper, error) {
	file, err := os.OpenFile(env.Storefile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("error during opening file: %s", err)
		return nil, err
	}
	log.Printf("sucsessifuly open file: %s", env.Storefile)
	return &Backuper{
		storage: storage,
		env:     env,
		enc:     enc,

		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

type Loader struct {
	storage *storage.MetricsStorage
	env     *env.ConfigServer
	encoder *encoder.Metrics

	file   *os.File
	reader *bufio.Reader
}

func NewLoader(storage *storage.MetricsStorage, env *env.ConfigServer, encoder *encoder.Metrics) (*Loader, error) {
	file, err := os.OpenFile(env.Storefile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("error during opening file: %s", err)
		return nil, err
	}
	log.Printf("sucsessifuly open file: %s", env.Storefile)
	return &Loader{
		storage: storage,
		env:     env,
		encoder: encoder,

		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (s *Backuper) Close() error {
	return s.file.Close()
}

func (s *Backuper) SaveToFile() error {
	for k, v := range s.storage.MetricsGauge {
		s.enc.ID = k
		s.enc.MType = "gauge"
		s.enc.Value = &v
		s.enc.Delta = nil

		data, err := s.enc.Marshall()
		if err != nil {
			log.Printf("error during marchalling: %s", err)
			return err
		}
		_, err = s.writer.Write(data)
		if err != nil {
			log.Printf("error duriong writing: %s", err)
			return err
		}
		err = s.writer.WriteByte('\n')
		if err != nil {
			return err
		}
	}
	for k, v := range s.storage.MetricsCounter {
		s.enc.ID = k
		s.enc.MType = "counter"
		s.enc.Value = nil
		s.enc.Delta = &v

		data, err := s.enc.Marshall()
		if err != nil {
			log.Printf("error during marchalling: %s", err)
			return err
		}
		_, err = s.writer.Write(data)
		if err != nil {
			log.Printf("error duriong writing: %s", err)
			return err
		}
		err = s.writer.WriteByte('\n')
		if err != nil {
			return err
		}
	}
	log.Printf("new backup created")
	return s.writer.Flush()
}

func (s *Backuper) Saver() error {
	tickerstore := time.NewTicker(s.env.StoreInterval)
	for {
		<-tickerstore.C
		err := s.SaveToFile()
		if err != nil {
			log.Printf("error during savihng: %s", err)
			return err
		}
	}
}

func (l *Loader) Close() error {
	return l.file.Close()
}

func (l *Loader) LoadFromFile() ([]byte, error) {

	for {
		data, err := l.reader.ReadBytes('\n')
		if err != nil {
			log.Printf("error during read file: %s", err)
			return nil, err
		}
		log.Printf("data %s", data)
		err = json.Unmarshal(data, l.encoder)
		if err != nil {
			log.Printf("error during unmarshalling: %s", err)
			return nil, err
		}
		switch l.encoder.MType {
		case "gauge":
			l.storage.InsertMetric(l.encoder.ID, *l.encoder.Value)
		case "counter":
			l.storage.CountCounterMetric(l.encoder.ID, *l.encoder.Delta)
		}
	}
}
