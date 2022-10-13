package service

import (
	"fmt"
	"time"

	"github.com/r2dtools/agent/pkg/logger"
)

const (
	DEFAULT_COLLECT_INTERVALL  = time.Minute
	DEFAULT_CLEAN_INTERVALL    = time.Hour * 24
	DEFAULT_MAX_COLLECT_PERIOD = 6 // month
)

var statCollectors []*StatCollector

type StatCollectorService struct {
	logger         logger.LoggerInterface
	statCollectors []*StatCollector
}

func (s *StatCollectorService) Run() error {
	for {
		for _, collector := range s.statCollectors {
			if err := collector.Collect(); err != nil {
				nErr := fmt.Errorf("could not collect data for '%s': %v", collector.Provider.GetCode(), err)
				s.logger.Error(nErr.Error())
			}
		}

		time.Sleep(DEFAULT_COLLECT_INTERVALL)
	}
}

func GetStatCollectorService(logger logger.LoggerInterface) (*StatCollectorService, error) {
	collectors, err := getStatCollectors(logger)
	if err != nil {
		return nil, err
	}
	return &StatCollectorService{statCollectors: collectors, logger: logger}, nil
}

type StatCleanerService struct {
	logger         logger.LoggerInterface
	statCollectors []*StatCollector
}

func GetStatCleanerService(logger logger.LoggerInterface) (*StatCleanerService, error) {
	collectors, err := getStatCollectors(logger)
	if err != nil {
		return nil, err
	}
	return &StatCleanerService{statCollectors: collectors, logger: logger}, nil
}

func (s *StatCleanerService) Run() error {
	for {
		toTime := time.Now().AddDate(0, -DEFAULT_MAX_COLLECT_PERIOD, 0).Unix()
		filter := StatProviderTimeFilter{
			FromTime: 0,
			ToTime:   int(toTime),
		}
		for _, collector := range s.statCollectors {
			if err := collector.Clean(&filter); err != nil {
				nErr := fmt.Errorf("could not clean up data for '%s': %v", collector.Provider.GetCode(), err)
				s.logger.Error(nErr.Error())
			}
		}

		time.Sleep(DEFAULT_CLEAN_INTERVALL)
	}
}

func getStatCollectors(logger logger.LoggerInterface) ([]*StatCollector, error) {
	if statCollectors != nil {
		return statCollectors, nil
	}

	var collectors []*StatCollector

	// cpu overall statistics
	cpuStatCollector, err := GetStatCollector(&OverallCPUStatPrivider{}, logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, cpuStatCollector)

	// cpu cores statistics
	cpuCoreStatCollectors, err := GetCoreCpuStatCollectors(logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, cpuCoreStatCollectors...)

	// virtual memory statistics
	virtualMemoryStatCollector, err := GetStatCollector(&VirtualMemoryStatPrivider{}, logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, virtualMemoryStatCollector)

	// swap statistics
	swapMemoryStatCollector, err := GetStatCollector(&SwapMemoryStatPrivider{}, logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, swapMemoryStatCollector)

	// disk usage statistics
	diskUsageStatCollector, err := GetDiskUsageStatCollector(logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, diskUsageStatCollector)

	// disk io statistics
	diskIOStatCollectors, err := GetDiskIOStatCollectors(logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, diskIOStatCollectors...)

	// network overall statistics
	networkStatCollector, err := GetStatCollector(&OverallNetworkStatProvider{}, logger)
	if err != nil {
		return nil, err
	}
	collectors = append(collectors, networkStatCollector)
	statCollectors = collectors

	return statCollectors, nil
}
