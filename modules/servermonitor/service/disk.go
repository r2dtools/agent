package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	diskService "github.com/r2dtools/agent/modules/servermonitor/service/disk"
	"github.com/shirou/gopsutil/disk"
)

// DiskUsageStatProvider retrieves statistics data for the disk usage
type DiskUsageStatProvider struct {
	Mapper *diskService.MountpointIDMapper
}

func (m *DiskUsageStatProvider) GetData() ([]string, error) {
	partitions, err := GetPartitions()
	if err != nil {
		return nil, err
	}

	usageData := make(map[int]string)
	for _, partition := range partitions {
		mountpointId, err := m.Mapper.GetMountpointID(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		usageData[mountpointId] = formatSpaceValue(usageStat.Used)
	}
	usageDataBytes, err := json.Marshal(usageData)
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, string(usageDataBytes))

	// time|{mountpointId: '', ....}
	return data, nil
}

func (m *DiskUsageStatProvider) GetCode() string {
	return DISK_USAGE_PROVIDER_CODE
}

func (m *DiskUsageStatProvider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 2 {
		return false
	}

	if filter == nil {
		return true
	}

	return filter.Check(data)
}

func (m *DiskUsageStatProvider) GetDiskInfo() (map[string]map[string]string, error) {
	partitions, err := GetPartitions()
	if err != nil {
		return nil, err
	}
	diskInfo := make(map[string]map[string]string)
	for _, partition := range partitions {
		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		mountpointId, err := m.Mapper.GetMountpointID(partition.Mountpoint)
		if err != nil {
			return nil, err
		}
		diskInfo[strconv.Itoa(mountpointId)] = map[string]string{
			"mountpoint":  partition.Mountpoint,
			"fstype":      usageStat.Fstype,
			"used":        formatSpaceValue(usageStat.Used),
			"free":        formatSpaceValue(usageStat.Free),
			"total":       formatSpaceValue(usageStat.Total),
			"usedPercent": fmt.Sprintf("%.2f", usageStat.UsedPercent),
		}
	}

	return diskInfo, nil
}

func formatSpaceValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}

func GetPartitions() ([]disk.PartitionStat, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("could not get partitions: %v", err)
	}

	var fPartitions []disk.PartitionStat
	for _, partition := range partitions {
		if strings.Contains(partition.Device, "/loop") || !strings.HasPrefix(partition.Device, "/dev") {
			continue
		}
		fPartitions = append(fPartitions, partition)
	}

	return fPartitions, nil
}
