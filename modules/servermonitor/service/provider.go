package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/modules/servermonitor/service/disk"
	"github.com/shirou/gopsutil/cpu"
)

const (
	OVERALL_CPU_PROVIDER_CODE     = "cpuoverall"
	CORE_CPU_PROVIDER_CODE        = "cpucore"
	VIRTUAL_MEMORY_PROVIDER_CODE  = "memoryvirtual"
	SWAP_MEMORY_PROVIDER_CODE     = "memoryswap"
	DISK_USAGE_PROVIDER_CODE      = "diskusage"
	DISK_IO_PROVIDER_CODE         = "diskio"
	OVERALL_NETWORK_PROVIDER_CODE = "networkoverall"

	DEFAULT_COLLECT_INTERVALL = time.Minute
)

type StatProvider interface {
	GetRecord() ([]string, error)
	GetCode() string
	CheckRecord([]string, StatProviderFilter) bool
	GetFieldsCount() int
	GetAverageRecord([][]string) []string
	GetEmptyRecordValue(index int) string
}

type BaseStatProvider struct{}

func (p *BaseStatProvider) GetEmptyRecordValue(index int) string {
	return ""
}

func (p *BaseStatProvider) getAverageRecord(records [][]string, fieldsCount int, formatFloat bool, getRecordEmptyValue func(int) string) []string {
	averageRecord := make([]string, fieldsCount)
	for i := 0; i < fieldsCount; i += 1 {
		var averageValue float64
		var recordsCount int
		for _, record := range records {
			if record[i] == getRecordEmptyValue(i) {
				continue
			}
			value, err := strconv.ParseFloat(record[i], 32)
			if err != nil {
				continue
			}
			recordsCount += 1
			averageValue += value
		}
		if recordsCount != 0 {
			averageValue = averageValue / float64(recordsCount)
			if formatFloat {
				averageRecord[i] = strconv.FormatFloat(averageValue, 'f', 2, 32)
			} else {
				averageRecord[i] = strconv.Itoa(int(averageValue))
			}
		}
	}

	return averageRecord
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
		providers = append(providers, &CoreCPUStatPrivider{CoreNum: i})
	}

	return providers, nil
}

func GetDiskUsageStatProvider() (StatProvider, error) {
	dataFolder := getDataFolder()
	if err := ensureFolderExists(dataFolder); err != nil {
		return nil, err
	}

	mounpointIdMapper, err := disk.GetMountpointIDMapper(dataFolder)
	if err != nil {
		return nil, err
	}

	return &DiskUsageStatProvider{Mapper: mounpointIdMapper}, nil
}
