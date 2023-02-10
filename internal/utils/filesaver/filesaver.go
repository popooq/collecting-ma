package filesaver

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/popooq/collectimg-ma/internal/server/config"
	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type Backuper struct {
	storage *storage.MetricsStorage
	cfg     *config.Config
	enc     *encoder.Encode

	file   *os.File
	writer *bufio.Writer
}

func NewBackuper(storage *storage.MetricsStorage, cfg *config.Config, enc *encoder.Encode) (*Backuper, error) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return nil, err
	}

	return &Backuper{
		storage: storage,
		cfg:     cfg,
		enc:     enc,

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
	s.file.Truncate(0)
	_ = s.writer.WriteByte('[')
	for k, v := range s.storage.MetricsGauge {
		s.fillEncGauge(k, v)

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
		s.fillEncCounter(k, d)

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
	return s.writer.Flush()
}

func (s *Backuper) GoFile() {
	tickerstore := time.NewTicker(s.cfg.StoreInterval)

	for {
		<-tickerstore.C

		err := s.SaveToFile()
		if err != nil {
			return
		}

	}
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

func (s *Backuper) fillEncGauge(k string, v float64) {
	s.enc.ID = k
	s.enc.MType = "gauge"
	s.enc.Value = &v
}

func (s *Backuper) fillEncCounter(k string, d int64) {
	s.enc.ID = k
	s.enc.MType = "counter"
	s.enc.Delta = &d
}

type Saver struct {
	cfg *config.Config

	file       *os.File
	readwriter *bufio.ReadWriter
}

func New(cfg *config.Config) (*Saver, error) {
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(file)

	return &Saver{
		cfg: cfg,

		file:       file,
		readwriter: bufio.NewReadWriter(reader, writer),
	}, nil
}

func (s *Saver) SaveMetric(metric *encoder.Encode) error {
	s.file.Truncate(0)
	_ = s.readwriter.WriteByte('[')

	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	_, err = s.readwriter.Write(data)
	if err != nil {
		return err
	}
	err = s.readwriter.WriteByte(',')
	if err != nil {
		return err
	}
	err = s.readwriter.WriteByte('\n')
	if err != nil {
		return err
	}
	_, _ = s.readwriter.WriteString("{}]")

	return s.readwriter.Flush()
}

func (s *Saver) SaveAllMetrics(metric encoder.Encode) error {
	s.file.Truncate(0)
	_ = s.readwriter.WriteByte('[')

	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	_, err = s.readwriter.Write(data)
	if err != nil {
		return err
	}
	err = s.readwriter.WriteByte(',')
	if err != nil {
		return err
	}
	err = s.readwriter.WriteByte('\n')
	if err != nil {
		return err
	}
	_, _ = s.readwriter.WriteString("{}]")

	return s.readwriter.Flush()
}

func (s *Saver) LoadMetrics() (metrics []encoder.Encode, err error) {
	var data []byte

	data, err = io.ReadAll(s.readwriter)
	if err != nil {
		log.Printf("erad err : %s", err)
		return nil, err
	}

	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Printf("marshal err : %s", err)
		return nil, err
	}

	return metrics, nil
}
