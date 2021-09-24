package handler

import (
	"errors"

	"github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agentintegration"
	"github.com/shirou/gopsutil/disk"
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

		diskProvider, ok := collector.Provider.(*service.DiskUsageStatProvider)
		if !ok {
			return errors.New("invalid disk statistics data provider type")
		}

		for _, row := range rows {
			diskUsageData = append(diskUsageData, getDiskUsageTimeLinePoint(row))
		}
		code := diskProvider.GetCode()
		responseData.TimeLineData[code] = diskUsageData

		usageData, err := collector.Provider.GetData()
		if err != nil {
			return nil
		}
		responseData.DiskInfo[code] = getDiskInfoItem(diskProvider.Partition, usageData)
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

// time|total|free|used|usedPercent
func getDiskInfoItem(partition disk.PartitionStat, row []string) map[string]string {
	item := make(map[string]string)
	item["total"] = row[1]
	item["free"] = row[2]
	item["used"] = row[3]
	item["usedPercent"] = row[4]
	item["fstype"] = partition.Fstype
	item["mountpoint"] = partition.Mountpoint
	item["device"] = partition.Device

	return item
}
