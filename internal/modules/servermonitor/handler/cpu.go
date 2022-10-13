package handler

import (
	"fmt"

	"github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

func LoadCpuTimeLineData(requestData *agentintegration.ServerMonitorStatisticsRequestData, logger logger.LoggerInterface) (*agentintegration.ServerMonitorStatisticsResponseData, error) {
	var responseData agentintegration.ServerMonitorStatisticsResponseData
	responseData.Data = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}
	if err := loadOverallCpuTimeLineData(&responseData, filter, logger); err != nil {
		return nil, err
	}
	if err := loadCoreCpuTimeLineData(&responseData, filter, logger); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadOverallCpuTimeLineData(responseData *agentintegration.ServerMonitorStatisticsResponseData, filter service.StatProviderFilter, logger logger.LoggerInterface) error {
	overallCpuStatCollector, err := service.GetStatCollector(&service.OverallCPUStatPrivider{}, logger)
	if err != nil {
		return err
	}

	rows, err := overallCpuStatCollector.Load(filter)
	if err != nil {
		return err
	}

	var overallCpuData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		overallCpuData = append(overallCpuData, getCpuTimeLinePoint(row))
	}
	responseData.Data["overall"] = overallCpuData

	return nil
}

func loadCoreCpuTimeLineData(responseData *agentintegration.ServerMonitorStatisticsResponseData, filter service.StatProviderFilter, logger logger.LoggerInterface) error {
	coreCpuStatCollectors, err := service.GetCoreCpuStatCollectors(logger)
	if err != nil {
		return err
	}

	for index, collector := range coreCpuStatCollectors {
		var coreCpuData []agentintegration.ServerMonitorTimeLinePoint
		rows, err := collector.Load(filter)
		if err != nil {
			return err
		}

		for _, row := range rows {
			coreCpuData = append(coreCpuData, getCpuTimeLinePoint(row))
		}
		responseData.Data[fmt.Sprintf("core%d", index+1)] = coreCpuData
	}

	return nil
}

func getCpuTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"system": row[1],
			"user":   row[2],
			"nice":   row[3],
			"idle":   row[4],
		},
	}
}
