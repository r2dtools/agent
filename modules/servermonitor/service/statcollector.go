package service

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/modules/servermonitor/service/disk"
	"github.com/shirou/gopsutil/cpu"
	"github.com/unknwon/com"
)

const (
	OVERALL_CPU_PROVIDER_CODE    = "cpuoverall"
	CORE_CPU_PROVIDER_CODE       = "cpucore"
	VIRTUAL_MEMORY_PROVIDER_CODE = "memoryvirtual"
	SWAP_MEMORY_PROVIDER_CODE    = "memoryswap"
	DISK_USAGE_PROVIDER_CODE     = "mountpoint"
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
	mu       *sync.RWMutex
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

	sc.mu.Lock()
	defer sc.mu.Unlock()

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

// GetCoreCpuStatProviders creates statistics providers for cpu cores
func GetCoreCpuStatProviders() ([]StatProvider, error) {
	cores, err := cpu.Counts(false)
	if err != nil {
		return nil, fmt.Errorf("could not create statisitcs providers for cpu cores: %v", err)
	}
	logger.Debug(fmt.Sprintf("count of cpu cores: %d", cores))

	var providers []StatProvider
	for i := 1; i <= cores; i++ {
		providers = append(providers, &CoreCPUStatPrivider{i})
	}

	return providers, nil
}

func GetCoreCpuStatCollectors() ([]*StatCollector, error) {
	providers, err := GetCoreCpuStatProviders()
	if err != nil {
		return nil, err
	}

	return GetStatCollectors(providers)
}

func GetStatCollectors(providers []StatProvider) ([]*StatCollector, error) {
	var collectors []*StatCollector
	for _, provider := range providers {
		collector, err := GetStatCollector(provider)
		if err != nil {
			logger.Debug(err.Error())
			continue
		}
		collectors = append(collectors, collector)
	}

	return collectors, nil
}

func GetDiskUsageStatCollectors() ([]*StatCollector, error) {
	providers, err := GetDiskUsageStatProviders()
	if err != nil {
		return nil, fmt.Errorf("could not create statistics providers for disk usage: %v", err)
	}

	return GetStatCollectors(providers)
}

func GetDiskUsageStatProviders() ([]StatProvider, error) {
	partitions, err := GetPartitions()
	if err != nil {
		return nil, err
	}

	dataFolder := getDataFolder()
	if err = ensureFolderExists(dataFolder); err != nil {
		return nil, err
	}

	var providers []StatProvider
	mounpointIdMapper, err := disk.GetMountpointIDMapper(dataFolder)
	if err != nil {
		return nil, err
	}
	for _, partition := range partitions {
		mountpointId, err := mounpointIdMapper.GetMountpointID(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		providers = append(providers, &DiskUsageStatProvider{Partition: partition, MountPointID: mountpointId})
	}

	return providers, nil
}

func GetStatCollector(provider StatProvider) (*StatCollector, error) {
	dataFolderPath := getDataFolder()
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

	return &StatCollector{&sync.RWMutex{}, provider, statFilePath}, nil
}

func getDataFolder() string {
	varDirPath := config.GetConfig().GetVarDirAbsPath()
	return filepath.Join(varDirPath, "modules", "servermonitor-module", "statistics")
}

func ensureFolderExists(folder string) error {
	if !com.IsDir(folder) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
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
