package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	diskService "github.com/r2dtools/agent/modules/servermonitor/service/disk"
	"github.com/shirou/gopsutil/disk"
	"github.com/unknwon/com"
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

func GetDiskDevices() ([]string, error) {
	partitions, err := GetPartitions()
	if err != nil {
		return nil, err
	}

	sdSubPartitionRegexp, err := regexp.Compile(`^/dev/(sd[a-z]+)(\d*)$`)
	if err != nil {
		return nil, err
	}
	nvmeSubPartitionRegexp, err := regexp.Compile(`^/dev/(nvme.+)(p\d*)$`)
	if err != nil {
		return nil, err
	}

	var diskDevices []string
	deviceRegexps := []*regexp.Regexp{sdSubPartitionRegexp, nvmeSubPartitionRegexp}
	subPartitionMap := make(map[string]string)

	for _, partition := range partitions {
		device := partition.Device
		for _, deviceRegexp := range deviceRegexps {
			groups := deviceRegexp.FindStringSubmatch(device)
			if len(groups) == 0 {
				continue
			}

			gDevice := groups[1]
			gSubPartition := groups[2]
			if gSubPartition != "" {
				subPartitionMap[gSubPartition] = gDevice
			} else if !com.IsSliceContainsStr(diskDevices, gDevice) {
				diskDevices = append(diskDevices, gDevice)
			}
		}
	}

	for subPartition, device := range subPartitionMap {
		if !com.IsSliceContainsStr(diskDevices, device) {
			diskDevices = append(diskDevices, device+subPartition)
		}
	}

	return diskDevices, nil
}

// DiskIOStatProvider retrieves statistics data for the disk IO
type DiskIOStatProvider struct {
	Device           string
	IOMeasureStorage *diskService.IOMeasureStorage
}

func (m *DiskIOStatProvider) GetData() ([]string, error) {
	ioStats, err := disk.IOCounters(m.Device)
	if err != nil {
		return nil, err
	}
	ioStat, ok := ioStats[m.Device]
	if !ok {
		return nil, nil
	}

	lastMeasure, err := m.IOMeasureStorage.GetLast(m.Device)
	if err != nil {
		return nil, err
	}
	if err = m.setLastMeasure(ioStat); err != nil {
		return nil, err
	}
	if lastMeasure == nil {
		return nil, nil
	}

	var data []string
	var readBytes, writeBytes uint64
	readTime := ioStat.ReadTime - lastMeasure.ReadTime
	writeTime := ioStat.WriteTime - lastMeasure.WriteTime
	ioTime := ioStat.IoTime - lastMeasure.IoTime

	if readTime != 0 {
		readBytes = (ioStat.ReadBytes - lastMeasure.ReadBytes) / readTime
	}
	if writeTime != 0 {
		writeBytes = (ioStat.WriteBytes - lastMeasure.WriteBytes) / writeTime
	}

	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatIOValue(ioStat.ReadCount-lastMeasure.ReadCount))
	data = append(data, formatIOValue(ioStat.WriteCount-lastMeasure.WriteCount))
	data = append(data, formatIOValue(ioStat.MergedReadCount-lastMeasure.MergedReadCount))
	data = append(data, formatIOValue(ioStat.MergedWriteCount-lastMeasure.MergedWriteCount))
	data = append(data, formatIOValue(readTime))
	data = append(data, formatIOValue(writeTime))
	data = append(data, formatIOValue(ioTime))
	data = append(data, formatIOValue(readBytes))
	data = append(data, formatIOValue(writeBytes))

	// time|ReadCount|WriteCount|MergedReadCount|MergedWriteCount|ReadTime|WriteTime|IoTime|ReadBytes|WriteBytes
	return data, nil
}

func (m *DiskIOStatProvider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 10 {
		return false
	}
	if filter == nil {
		return true
	}
	return filter.Check(data)
}

func (m *DiskIOStatProvider) GetCode() string {
	return DISK_IO_PROVIDER_CODE + m.Device
}

func (m *DiskIOStatProvider) setLastMeasure(ioStat disk.IOCountersStat) error {
	measure := &diskService.IOMeasure{
		ReadCount:        ioStat.ReadCount,
		WriteCount:       ioStat.WriteCount,
		MergedReadCount:  ioStat.MergedReadCount,
		MergedWriteCount: ioStat.MergedWriteCount,
		ReadTime:         ioStat.ReadTime,
		WriteTime:        ioStat.WriteTime,
		IoTime:           ioStat.IoTime,
		ReadBytes:        ioStat.ReadBytes,
		WriteBytes:       ioStat.WriteBytes,
	}
	if err := m.IOMeasureStorage.SetLast(m.Device, measure); err != nil {
		return err
	}
	return nil
}

func formatSpaceValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}

func formatIOValue(value uint64) string {
	return strconv.FormatUint(value, 10)
}
