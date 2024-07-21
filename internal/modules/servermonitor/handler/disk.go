package handler

import (
	"encoding/json"
	"errors"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

func LoadDiskUsageTimeLineData(
	requestData *agentintegration.ServerMonitorStatisticsRequestData,
	config *config.Config,
	logger logger.Logger,
) (*agentintegration.ServerMonitorDiskResponseData, error) {
	var responseData agentintegration.ServerMonitorDiskResponseData
	responseData.DiskUsageTimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.DiskIOTimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.DiskInfo = make(map[string]map[string]string)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}

	if err := loadDiskUsageData(&responseData, filter, config, logger); err != nil {
		return nil, err
	}

	if err := loadDiskIOData(&responseData, filter, config, logger); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadDiskIOData(
	responseData *agentintegration.ServerMonitorDiskResponseData,
	filter service.StatProviderFilter,
	config *config.Config,
	logger logger.Logger,
) error {
	diskIOStatCollectors, err := service.GetDiskIOStatCollectors(config, logger)
	if err != nil {
		return err
	}

	for _, collector := range diskIOStatCollectors {
		var diskIoData []agentintegration.ServerMonitorTimeLinePoint
		rows, err := collector.Load(filter)
		if err != nil {
			return err
		}

		for _, row := range rows {
			diskIoData = append(diskIoData, getDiskIOTimeLineData(row))
		}
		provider, ok := collector.Provider.(*service.DiskIOStatProvider)
		if !ok {
			return errors.New("invalid type of disk io statistics provider")
		}
		responseData.DiskIOTimeLineData[provider.Device] = diskIoData
	}

	return nil
}

func loadDiskUsageData(
	responseData *agentintegration.ServerMonitorDiskResponseData,
	filter service.StatProviderFilter,
	config *config.Config,
	logger logger.Logger,
) error {
	diskUsageStatCollector, err := service.GetDiskUsageStatCollector(config, logger)
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
	responseData.DiskUsageTimeLineData[diskUsageStatProvider.GetCode()] = diskUsageData
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

// row: time|ReadCount|WriteCount|MergedReadCount|MergedWriteCount|ReadTime|WriteTime|IoTime|ReadBytes|WriteBytes
func getDiskIOTimeLineData(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"readcount":        row[1],
			"writecount":       row[2],
			"mergedreadcount":  row[3],
			"mergedwritecount": row[4],
			"readtime":         row[5],
			"writetime":        row[6],
			"iotime":           row[7],
			"readbytes":        row[8],
			"writebytes":       row[9],
		},
	}
}
