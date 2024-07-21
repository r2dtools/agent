package handler

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

func LoadMemoryTimeLineData(
	requestData *agentintegration.ServerMonitorStatisticsRequestData,
	config *config.Config,
	logger logger.Logger,
) (*agentintegration.ServerMonitorStatisticsResponseData, error) {
	var responseData agentintegration.ServerMonitorStatisticsResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}

	if err := loadVirtualMemoryTimeLineData(&responseData, filter, config, logger); err != nil {
		return nil, err
	}

	if err := loadSwapMemoryTimeLineData(&responseData, filter, config, logger); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadVirtualMemoryTimeLineData(
	responseData *agentintegration.ServerMonitorStatisticsResponseData,
	filter service.StatProviderFilter,
	config *config.Config,
	logger logger.Logger,
) error {
	virtualMemoryStatCollector, err := service.GetStatCollector(&service.VirtualMemoryStatPrivider{}, config, logger)
	if err != nil {
		return nil
	}

	rows, err := virtualMemoryStatCollector.Load(filter)
	if err != nil {
		return nil
	}

	var memoryData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		memoryData = append(memoryData, getVirtualMemoryTimeLinePoint(row))
	}
	responseData.Data["virtual"] = memoryData

	return nil
}

func loadSwapMemoryTimeLineData(
	responseData *agentintegration.ServerMonitorStatisticsResponseData,
	filter service.StatProviderFilter,
	config *config.Config,
	logger logger.Logger,
) error {
	swapMemoryStatCollector, err := service.GetStatCollector(&service.SwapMemoryStatPrivider{}, config, logger)
	if err != nil {
		return nil
	}

	rows, err := swapMemoryStatCollector.Load(filter)
	if err != nil {
		return nil
	}

	var memoryData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		memoryData = append(memoryData, getSwapMemoryTimeLinePoint(row))
	}
	responseData.Data["swap"] = memoryData

	return nil
}

// time|total|available|free|used|active|inactive|cached|buffered
func getVirtualMemoryTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"total":     row[1],
			"available": row[2],
			"free":      row[3],
			"used":      row[4],
			"active":    row[5],
			"inactive":  row[6],
			"cached":    row[7],
			"buffered":  row[8],
		},
	}
}

// time|total|used|free
func getSwapMemoryTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"total": row[1],
			"used":  row[2],
			"free":  row[3],
		},
	}
}
