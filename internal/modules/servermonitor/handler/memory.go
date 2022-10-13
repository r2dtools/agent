package handler

import (
	"github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

func LoadMemoryTimeLineData(requestData *agentintegration.ServerMonitorStatisticsRequestData, logger logger.LoggerInterface) (*agentintegration.ServerMonitorStatisticsResponseData, error) {
	var responseData agentintegration.ServerMonitorStatisticsResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}

	if err := loadVirtualMemoryTimeLineData(&responseData, filter, logger); err != nil {
		return nil, err
	}
	if err := loadSwapMemoryTimeLineData(&responseData, filter, logger); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadVirtualMemoryTimeLineData(responseData *agentintegration.ServerMonitorStatisticsResponseData, filter service.StatProviderFilter, logger logger.LoggerInterface) error {
	virtualMemoryStatCollector, err := service.GetStatCollector(&service.VirtualMemoryStatPrivider{}, logger)
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

func loadSwapMemoryTimeLineData(responseData *agentintegration.ServerMonitorStatisticsResponseData, filter service.StatProviderFilter, logger logger.LoggerInterface) error {
	swapMemoryStatCollector, err := service.GetStatCollector(&service.SwapMemoryStatPrivider{}, logger)
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
