package service

import (
	"fmt"
	"time"

	"github.com/r2dtools/agent/logger"
)

type Service struct {
	statCollectors []*StatCollector
}

func GetService() (*Service, error) {
	var collectors []*StatCollector

	cpuStatCollector, err := GetStatCollector(&OverallCPUStatPrivider{})
	if err != nil {
		return nil, err
	}

	collectors = append(collectors, cpuStatCollector)

	return &Service{collectors}, nil
}

func (s *Service) Run() error {
	for {
		for _, collector := range s.statCollectors {
			if err := collector.Collect(); err != nil {
				nErr := fmt.Errorf("could not collect data for '%s': %v", collector.Provider.GetCode(), err)
				logger.Error(nErr.Error())
			}
		}

		time.Sleep(time.Minute)
	}
}
