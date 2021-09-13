package service

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/logger"
	"github.com/unknwon/com"
)

const (
	OVERALL_CPU_PROVIDER_CODE = "overall_cpu"
)

type StatProvider interface {
	GetData() ([]string, error)
	GetCode() string
	CheckData([]string, StatProviderFilter) bool
}

type StatProviderFilter interface {
	Check(row []string) bool
}

type StatCollector struct {
	Provider StatProvider
	FilePath string
}

func (sc *StatCollector) Collect() error {
	data, err := sc.Provider.GetData()
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	file, err := os.OpenFile(sc.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '|'
	if err = writer.Write(data); err != nil {
		return fmt.Errorf("could not write statistics data for '%s': %v", sc.Provider.GetCode(), err)
	}
	writer.Flush()

	if err = writer.Error(); err != nil {
		return fmt.Errorf("could not write statistics data for '%s': %v", sc.Provider.GetCode(), err)
	}

	return nil
}

func (sc *StatCollector) Load(filter StatProviderFilter) ([][]string, error) {
	file, err := os.Open(sc.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// buffer size 100kb
	bReader := bufio.NewReaderSize(file, 102400)
	reader := csv.NewReader(bReader)
	reader.Comma = '|'
	var data [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Debug(fmt.Sprintf("could not read record from '%s' collector: %v", sc.Provider.GetCode(), err))
			continue
		}

		if !sc.Provider.CheckData(record, filter) {
			continue
		}
		data = append(data, record)
	}

	return data, nil
}

func GetStatCollector(provider StatProvider) (*StatCollector, error) {
	dataFolderPath := getDataFolder()
	if !com.IsDir(dataFolderPath) {
		if err := os.MkdirAll(dataFolderPath, 0755); err != nil {
			return nil, fmt.Errorf("could not create statistics collector '%s': %v", provider.GetCode(), err)
		}
	}

	statFilePath := filepath.Join(dataFolderPath, provider.GetCode())
	if !com.IsFile(statFilePath) {
		_, err := os.Create(statFilePath)

		if err != nil {
			return nil, fmt.Errorf("could not create statistics collector '%s': %v", provider.GetCode(), err)
		}
	}

	return &StatCollector{provider, statFilePath}, nil
}

func getDataFolder() string {
	varDirPath := config.GetConfig().GetVarDirAbsPath()
	return filepath.Join(varDirPath, "modules", "servermonitor-module", "statistics")
}

type StatProviderTimeFilter struct {
	FromTime, ToTime int
}

func (f *StatProviderTimeFilter) Check(row []string) bool {
	if len(row) == 0 {
		return false
	}

	time, err := strconv.Atoi(row[0])
	if err != nil {
		return false
	}

	if f.FromTime > 0 && time < f.FromTime {
		return false
	}

	if f.ToTime > 0 && time > f.ToTime {
		return false
	}

	return true
}
