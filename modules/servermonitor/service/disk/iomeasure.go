package disk

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/unknwon/com"
)

type IOMeasureStorage struct {
	filePath string
}

//|ReadCount|WriteCount|MergedReadCount|MergedWriteCount|ReadTime|WriteTime|IoTime|ReadBytes|WriteBytes
type IOMeasure struct {
	ReadCount,
	WriteCount,
	MergedReadCount,
	MergedWriteCount,
	ReadTime,
	WriteTime,
	IoTime,
	ReadBytes,
	WriteBytes uint64
}

func GetIOMeasure(dataFolder string) (*IOMeasureStorage, error) {
	filePath := filepath.Join(dataFolder, "diskiomeasure")
	if !com.IsFile(filePath) {
		if _, err := os.Create(filePath); err != nil {
			return nil, err
		}
	}
	storage := IOMeasureStorage{filePath}

	return &storage, nil
}

func (m *IOMeasureStorage) GetLast(device string) (*IOMeasure, error) {
	file, err := os.OpenFile(m.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := m.getData(file)
	measure, ok := data[device]
	if !ok {
		return nil, nil
	}

	return measure, nil
}

func (m *IOMeasureStorage) SetLast(device string, measure *IOMeasure) error {
	file, err := os.OpenFile(m.filePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data := m.getData(file)
	if err = file.Truncate(0); err != nil {
		return err
	}
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	data[device] = measure
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (m *IOMeasureStorage) getData(file *os.File) map[string]*IOMeasure {
	data := make(map[string]*IOMeasure)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return data
	}

	return data
}
