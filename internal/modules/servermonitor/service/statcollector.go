package service

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/unknwon/com"
)

type StatCollector struct {
	mu       *sync.RWMutex
	Provider StatProvider
	FilePath string
	Logger   logger.Logger
	Config   *config.Config
}

func (sc *StatCollector) Collect() error {
	data, err := sc.Provider.GetRecord()
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	file, err := os.OpenFile(sc.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	cTime := strconv.FormatInt(time.Now().Unix(), 10)
	fullData := []string{cTime}
	fullData = append(fullData, data...)

	writer := csv.NewWriter(file)
	writer.Comma = '|'
	if err = writer.Write(fullData); err != nil {
		return fmt.Errorf("could not write statistics data for '%s': %v", sc.Provider.GetCode(), err)
	}
	writer.Flush()

	if err = writer.Error(); err != nil {
		return fmt.Errorf("could not write statistics data for '%s': %v", sc.Provider.GetCode(), err)
	}

	return nil
}

func (sc *StatCollector) Load(filter StatProviderFilter) ([][]string, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	file, err := os.Open(sc.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// buffer size 100kb
	bReader := bufio.NewReaderSize(file, 102400)
	reader := csv.NewReader(bReader)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1
	data := make([][]string, 0)
	var previousRecord []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			sc.Logger.Debug("could not read record from '%s' collector: %v", sc.Provider.GetCode(), err)
			continue
		}

		if !sc.Provider.CheckRecord(record, filter) {
			continue
		}

		data = append(data, sc.getExtendedRecords(previousRecord, record, DEFAULT_COLLECT_INTERVALL)...)
		previousRecord = record
	}
	data = sc.averageRecords(data, filter.GetFromTime(), filter.GetToTime())

	return data, nil
}

func (sc *StatCollector) Clean(filter StatProviderFilter) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	file, err := os.Open(sc.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1
	recordsToRemove := [][]string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			sc.Logger.Debug("could not read record from '%s' collector when data cleaning: %v", sc.Provider.GetCode(), err)
			continue
		}

		if sc.Provider.CheckRecord(record, filter) {
			recordsToRemove = append(recordsToRemove, record)
		} else {
			break
		}
	}

	if len(recordsToRemove) == 0 {
		return nil
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	var pos int64
	linesCount := 0
	bReader := bufio.NewReader(file)

	for {
		bRecord, err := bReader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		pos += int64(len(bRecord))
		linesCount++
		if linesCount >= len(recordsToRemove) {
			break
		}
	}

	_, err = file.Seek(pos, 0)
	if err != nil {
		return err
	}

	nFilePath := sc.FilePath + ".tmp"
	nFile, err := os.Create(nFilePath)
	if err != nil {
		return err
	}
	defer nFile.Close()

	_, err = io.Copy(nFile, file)
	if err != nil {
		return err
	}
	if err := nFile.Sync(); err != nil {
		return err
	}

	if err := os.Rename(nFilePath, sc.FilePath); err != nil {
		return err
	}

	return nil
}

func (sc *StatCollector) averageRecords(records [][]string, fromTime, toTime int) [][]string {
	if len(records) == 0 {
		return records
	}
	firstRecord := records[0]
	recordTime, _, err := sc.parseRecord(firstRecord)
	if err != nil {
		sc.Logger.Debug("could not parse first record of data for the provider '%s': %v", sc.Provider.GetCode(), err)
		return records
	}

	minTime := fromTime
	if fromTime < recordTime {
		minTime = recordTime
	}

	maxTime := toTime
	currentTime := time.Now().Unix()
	if toTime > int(currentTime) {
		maxTime = int(currentTime)
	}

	averageCount := sc.getAverageCount(minTime, maxTime)
	if averageCount <= 0 {
		return records
	}

	var averagedRecords [][]string
	for i := 0; i < len(records); i += averageCount {
		end := i + averageCount
		if end > len(records) {
			end = len(records)
		}
		chunk := records[i:end]
		var averageChunk [][]string
		var times []int
		for _, chunkItem := range chunk {
			time, values, err := sc.parseRecord(chunkItem)
			if err != nil {
				continue
			}
			times = append(times, time)
			averageChunk = append(averageChunk, values)
		}

		if len(times) == 0 {
			continue
		}
		lTime := times[len(times)-1]
		averageRecord := sc.Provider.GetAverageRecord(averageChunk)
		averageRecord = append([]string{strconv.Itoa(lTime)}, averageRecord...)
		averagedRecords = append(averagedRecords, averageRecord)
	}

	return averagedRecords
}

func (sc *StatCollector) getAverageCount(minTime, maxTime int) int {
	days := (float64(maxTime) - float64(minTime)) / (60 * 60 * 24) //days count
	maxDays := math.Ceil(days)
	if maxDays <= 1 {
		return -1
	}
	count := math.Ceil(maxDays / 7)

	return int(count) * 5
}

func (sc *StatCollector) getFullRecord(fieldsCount int, record []string) []string {
	nRecord := sc.getEmptyRecord(-1, fieldsCount)
	iterateCount := len(record)
	if fieldsCount == iterateCount {
		return record
	}
	if iterateCount > fieldsCount {
		iterateCount = fieldsCount
	}

	for i := 0; i < iterateCount; i += 1 {
		nRecord[i] = record[i]
	}

	return nRecord
}

func (sc *StatCollector) getEmptyRecord(time, fieldsCount int) []string {
	nRecord := make([]string, fieldsCount)
	for i := range nRecord {
		if i == 0 {
			nRecord[i] = ""
		} else {
			nRecord[i] = sc.Provider.GetEmptyRecordValue(i)
		}
	}
	if time > 0 {
		nRecord[0] = strconv.Itoa(time)
	}

	return nRecord
}

func (sc *StatCollector) getExtendedRecords(previousRecord, record []string, interval time.Duration) [][]string {
	fieldsCount := sc.Provider.GetFieldsCount() + 1
	eRecords := [][]string{sc.getFullRecord(fieldsCount, record)}

	if previousRecord == nil {
		return eRecords
	}

	lTime, _, err := sc.parseRecord(previousRecord)
	if err != nil {
		return eRecords
	}

	time, _, err := sc.parseRecord(record)
	if err != nil {
		return eRecords
	}

	timeDiff := time - lTime
	intervalSeconds := interval.Seconds()

	var addRecords [][]string
	for timeDiff > int(math.Ceil(intervalSeconds))+int(interval.Seconds()/2) {
		time := lTime + int(intervalSeconds)
		fullRecord := sc.getEmptyRecord(time, fieldsCount)
		addRecords = append(addRecords, fullRecord)
		timeDiff = timeDiff - int(intervalSeconds)
		lTime = time
	}
	eRecords = append(addRecords, eRecords...)

	return eRecords
}

func (sc *StatCollector) parseRecord(record []string) (int, []string, error) {
	if len(record) == 0 {
		return 0, nil, errors.New("record is empty")
	}
	time, err := strconv.Atoi(record[0])
	if err != nil {
		return 0, nil, err
	}
	return time, record[1:], nil
}

func GetCoreCpuStatCollectors(config *config.Config, logger logger.Logger) ([]*StatCollector, error) {
	providers, err := GetCoreCpuStatProviders(config, logger)
	if err != nil {
		return nil, err
	}

	return GetStatCollectors(providers, config, logger)
}

func GetStatCollectors(providers []StatProvider, config *config.Config, logger logger.Logger) ([]*StatCollector, error) {
	var collectors []*StatCollector
	for _, provider := range providers {
		collector, err := GetStatCollector(provider, config, logger)
		if err != nil {
			logger.Debug(err.Error())
			continue
		}
		collectors = append(collectors, collector)
	}

	return collectors, nil
}

func GetDiskUsageStatCollector(config *config.Config, logger logger.Logger) (*StatCollector, error) {
	provider, err := GetDiskUsageStatProvider(config)

	if err != nil {
		return nil, fmt.Errorf("could not create statistics provider for disk usage: %v", err)
	}

	return GetStatCollector(provider, config, logger)
}

func GetDiskIOStatCollectors(config *config.Config, logger logger.Logger) ([]*StatCollector, error) {
	dataFolder := getDataFolder(config)
	if err := ensureFolderExists(dataFolder); err != nil {
		return nil, err
	}

	devices, err := GetDiskDevices()
	if err != nil {
		return nil, err
	}

	var providers []StatProvider
	for _, device := range devices {
		providers = append(providers, &DiskIOStatProvider{Device: device})
	}

	return GetStatCollectors(providers, config, logger)
}

func GetStatCollector(provider StatProvider, config *config.Config, logger logger.Logger) (*StatCollector, error) {
	dataFolderPath := getDataFolder(config)

	if err := ensureFolderExists(dataFolderPath); err != nil {
		return nil, fmt.Errorf("could not create statistics collector '%s': %v", provider.GetCode(), err)
	}

	statFilePath := filepath.Join(dataFolderPath, provider.GetCode())

	if !com.IsFile(statFilePath) {
		_, err := os.Create(statFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not create statistics collector '%s': %v", provider.GetCode(), err)
		}
	}

	return &StatCollector{&sync.RWMutex{}, provider, statFilePath, logger, config}, nil
}

func getDataFolder(config *config.Config) string {
	return filepath.Join(config.GetModuleVarAbsDir("servermonitor"), "statistics")
}

func ensureFolderExists(folder string) error {
	if !com.IsDir(folder) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}
