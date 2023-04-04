package filesaver

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

// Saver структура хранит иенформацию о сейвере
type Saver struct {
	file       *os.File
	readwriter *bufio.ReadWriter
}

// New функция создает новый Saver
func New(storefile string) (*Saver, error) {
	file, err := os.OpenFile(storefile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(file)

	return &Saver{
		file:       file,
		readwriter: bufio.NewReadWriter(reader, writer),
	}, nil
}

// SaveMetric метод сохраняет метрики в файл
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

// LoadMetric метод загружает метрики
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

// KeeperCheck метод-заглушка для реализации интерфейса Keeper
func (s *Saver) KeeperCheck() error {
	return nil
}
