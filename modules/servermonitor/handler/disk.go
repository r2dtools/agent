package handler

import (
	"encoding/json"
	"errors"

	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
)

func LoadDiskUsageTimeLineData(requestData *agentintegration.ServerMonitorTimeLineRequestData) (*agentintegration.ServerMonitorDiskResponseData, error) {
	var responseData agentintegration.ServerMonitorDiskResponseData
	responseData.TimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.DiskInfo = make(map[string]map[string]string)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}
	if err := loadDiskUsageData(&responseData, filter); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadDiskUsageData(responseData *agentintegration.ServerMonitorDiskResponseData, filter service.StatProviderFilter) error {
	diskUsageStatCollector, err := service.GetDiskUsageStatCollector()
	if err != nil {
		return err
	}

	var diskUsageData []agentintegration.ServerMonitorTimeLinePoint
	rows, err := diskUsageStatCollector.Load(filter)
	if err != nil {
		return err
	}

	for _, row := range rows {
		diskUsageTimeLinePoint, err := getDiskUsageTimeLinePoint(row)
		if err == nil {
			diskUsageData = append(diskUsageData, diskUsageTimeLinePoint)
		}
	}

	provider := diskUsageStatCollector.Provider
	diskUsageStatProvider, ok := provider.(*service.DiskUsageStatProvider)
	if !ok {
		return errors.New("invalid type of disk usage statistics provider")
	}

	diskInfo, err := diskUsageStatProvider.GetDiskInfo()
	if err != nil {
		return err
	}
	responseData.TimeLineData[diskUsageStatProvider.GetCode()] = diskUsageData
	responseData.DiskInfo = diskInfo

	return nil
}

func getDiskUsageTimeLinePoint(row []string) (agentintegration.ServerMonitorTimeLinePoint, error) {
	usageData := make(map[string]string)
	var timeLinePoint agentintegration.ServerMonitorTimeLinePoint
	if err := json.Unmarshal([]byte(row[1]), &usageData); err != nil {
		return timeLinePoint, err
	}
	timeLinePoint = agentintegration.ServerMonitorTimeLinePoint{
		Time:  row[0],
		Value: usageData,
	}

	return timeLinePoint, nil
}
