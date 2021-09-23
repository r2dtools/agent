package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// DiskStatPrivider retrieves statistics data for the disk usage
type DiskStatPrivider struct {
	Partition *disk.PartitionStat
}

func (m *DiskStatPrivider) GetData() ([]string, error) {
	usageStat, err := disk.Usage(m.Partition.Mountpoint)
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatSpaceValue(usageStat.Total))
	data = append(data, formatSpaceValue(usageStat.Free))
	data = append(data, formatSpaceValue(usageStat.Used))
	data = append(data, fmt.Sprintf("%.2f", usageStat.UsedPercent))

	return data, nil
}

func (m *DiskStatPrivider) GetCode() string {
	return DISK_PROVIDER_CODE + strings.ReplaceAll(m.Partition.Device, "/", "")
}

func (m *DiskStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if filter == nil {
		return true
	}

	return filter.Check(data)
}

func formatSpaceValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}
