package handler

import (
	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
)

func LoadDiskUsageTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorDiskResponseData, error) {
	var responseData agentintegration.ServerMonitorDiskResponseData
	responseData.TimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}
	if err := loadDiskUsageTimeLineData(&responseData, filter); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadDiskUsageTimeLineData(responseData *agentintegration.ServerMonitorDiskResponseData, filter service.StatProviderFilter) error {
	diskUsageStatCollectors, err := service.GetDiskUsageStatCollectors()
	if err != nil {
		return err
	}

	for _, collector := range diskUsageStatCollectors {
		var diskUsageData []agentintegration.ServerMonitorTimeLinePoint
		rows, err := collector.Load(filter)
		if err != nil {
			return err
		}

		for _, row := range rows {
			diskUsageData = append(diskUsageData, getDiskUsageTimeLinePoint(row))
		}
		responseData.TimeLineData[collector.Provider.GetCode()] = diskUsageData
	}

	return nil
}

func getDiskUsageTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"used":        row[3],
			"usedPercert": row[4],
		},
	}
}
