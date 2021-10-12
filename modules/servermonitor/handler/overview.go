package handler

import (
	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
)

func LoadOverviewStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorOverviewResponseData, error) {
	var responseData agentintegration.ServerMonitorOverviewResponseData
	if err := loadOverviewCpuData(&responseData); err != nil {
		return nil, err
	}
	if err := loadOverviewDiskData(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func loadOverviewCpuData(responseData *agentintegration.ServerMonitorOverviewResponseData) error {
	overallCpuStatCollector, err := service.GetStatCollector(&service.OverallCPUStatPrivider{})
	if err != nil {
		return err
	}

	rows, err := overallCpuStatCollector.Load(nil)
	if err != nil {
		return err
	}

	var cpuUsageData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		cpuUsageData = append(cpuUsageData, getOverviewCpuTimeLinePoint(row))
	}
	responseData.CpuTimeLineData = cpuUsageData

	return nil
}

func loadOverviewDiskData(responseData *agentintegration.ServerMonitorOverviewResponseData) error {
	diskUsageStatCollector, err := service.GetDiskUsageStatCollector()
	if err != nil {
		return err
	}

	var diskUsageData []agentintegration.ServerMonitorTimeLinePoint
	rows, err := diskUsageStatCollector.Load(nil)
	if err != nil {
		return err
	}

	for _, row := range rows {
		diskUsageTimeLinePoint, err := getDiskUsageTimeLinePoint(row)
		if err == nil {
			diskUsageData = append(diskUsageData, diskUsageTimeLinePoint)
		}
	}
	responseData.DiskTimeLineData = diskUsageData

	return nil
}

func getOverviewCpuTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	var usage string
	if len(row) < 6 {
		usage = "0"
	} else {
		usage = row[5]
	}

	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"usage": usage,
		},
	}
}
