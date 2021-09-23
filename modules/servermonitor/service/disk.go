package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// DiskUsageStatPrivider retrieves statistics data for the disk usage
type DiskUsageStatPrivider struct {
	Mountpoint   string
	MountPointID int
}

func (m *DiskUsageStatPrivider) GetData() ([]string, error) {
	usageStat, err := disk.Usage(m.Mountpoint)
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatSpaceValue(usageStat.Total))
	data = append(data, formatSpaceValue(usageStat.Free))
	data = append(data, formatSpaceValue(usageStat.Used))
	data = append(data, fmt.Sprintf("%.2f", usageStat.UsedPercent))

	// time|total|free|used|usedPercent
	return data, nil
}

func (m *DiskUsageStatPrivider) GetCode() string {
	return fmt.Sprintf("%s%d", DISK_USAGE_PROVIDER_CODE, m.MountPointID)
}

func (m *DiskUsageStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 5 {
		return false
	}

	if filter == nil {
		return true
	}

	return filter.Check(data)
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
