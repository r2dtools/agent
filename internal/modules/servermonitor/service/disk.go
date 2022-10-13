package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	diskService "github.com/r2dtools/agent/internal/modules/servermonitor/service/disk"
	"github.com/shirou/gopsutil/disk"
	"github.com/unknwon/com"
)

type IOMeasure struct {
	ReadCount,
	WriteCount,
	MergedReadCount,
	MergedWriteCount,
	ReadTime,
	WriteTime,
	IoTime,
	ReadBytes,
	WriteBytes uint64
}

func init() {
	lastIOMeasure = make(map[string]IOMeasure)
}

var lastIOMeasure map[string]IOMeasure

// DiskUsageStatProvider retrieves statistics data for the disk usage
type DiskUsageStatProvider struct {
	BaseStatProvider
	Mapper *diskService.MountpointIDMapper
}

func (m *DiskUsageStatProvider) GetRecord() ([]string, error) {
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
	data = append(data, string(usageDataBytes))

	// {mountpointId: '', ....}
	return data, nil
}

func (m *DiskUsageStatProvider) GetEmptyRecordValue(index int) string {
	return "{}"
}

func (m *DiskUsageStatProvider) GetAverageRecord(records [][]string) []string {
	averageRecordSum := make(map[int]int)
	averageRecordCount := make(map[int]int)
	averageRecord := make(map[int]string)
	for _, record := range records {
		if len(record) == 0 || record[0] == m.GetEmptyRecordValue(0) {
			continue
		}

		recordData := make(map[int]string)
		if err := json.Unmarshal([]byte(record[0]), &recordData); err != nil {
			continue
		}
		for mountpoint, sValue := range recordData {
			value, err := strconv.Atoi(sValue)
			if err != nil {
				continue
			}
			averageRecordSum[mountpoint] += value
			averageRecordCount[mountpoint] += 1
		}
	}

	for mountpoint, sum := range averageRecordSum {
		averageRecord[mountpoint] = strconv.Itoa(sum / averageRecordCount[mountpoint])
	}
	data := make([]string, 1)
	averageRecordBytes, err := json.Marshal(averageRecord)
	if err != nil {
		data[0] = m.GetEmptyRecordValue(0)
	} else {
		data[0] = string(averageRecordBytes)
	}

	return data
}

func (m *DiskUsageStatProvider) GetFieldsCount() int {
	return 1
}

func (m *DiskUsageStatProvider) GetCode() string {
	return DISK_USAGE_PROVIDER_CODE
}

func (m *DiskUsageStatProvider) CheckRecord(data []string, filter StatProviderFilter) bool {
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

	return getDiskDevicesFromPartitions(partitions)
}

func getDiskDevicesFromPartitions(partitions []disk.PartitionStat) ([]string, error) {
	sdSubPartitionRegexp, err := regexp.Compile(`^/dev/(sd[a-z]+)(\d*)$`)
	if err != nil {
		return nil, err
	}
	nvmeSubPartitionRegexp, err := regexp.Compile(`^/dev/(nvme.+)(p\d*)$`)
	if err != nil {
		return nil, err
	}
	dmPartitionRegexp, err := regexp.Compile(`^/dev/(dm\-\d+)$`)
	if err != nil {
		return nil, err
	}

	var diskDevices []string
	deviceRegexps := []*regexp.Regexp{dmPartitionRegexp, sdSubPartitionRegexp, nvmeSubPartitionRegexp}
	subPartitionMap := make(map[string]string)

	for _, partition := range partitions {
		device := partition.Device
		for _, deviceRegexp := range deviceRegexps {
			groups := deviceRegexp.FindStringSubmatch(device)
			if len(groups) == 0 {
				continue
			}

			gDevice := groups[1]
			var gSubPartition string
			if len(groups) > 2 {
				gSubPartition = groups[2]
			}

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
	BaseStatProvider
	Device string
}

func (m *DiskIOStatProvider) GetRecord() ([]string, error) {
	ioStats, err := disk.IOCounters(m.Device)
	if err != nil {
		return nil, err
	}
	ioStat, ok := ioStats[m.Device]
	if !ok {
		return nil, nil
	}

	lastMeasure, ok := lastIOMeasure[m.Device]
	if !ok {
		lastIOMeasure[m.Device] = getLastIOMeasure(ioStat)
		return nil, nil
	}

	var data []string
	var readBytes, writeBytes uint64
	readTime := m.getDiff(ioStat.ReadTime, lastMeasure.ReadTime)
	writeTime := m.getDiff(ioStat.WriteTime, lastMeasure.WriteTime)
	ioTime := m.getDiff(ioStat.IoTime, lastMeasure.IoTime)

	if readTime != 0 {
		readBytes = m.getDiff(ioStat.ReadBytes, lastMeasure.ReadBytes) / readTime
	}
	if writeTime != 0 {
		writeBytes = m.getDiff(ioStat.WriteBytes, lastMeasure.WriteBytes) / writeTime
	}

	data = append(data, formatIOValue(m.getDiff(ioStat.ReadCount, lastMeasure.ReadCount)))
	data = append(data, formatIOValue(m.getDiff(ioStat.WriteCount, lastMeasure.WriteCount)))
	data = append(data, formatIOValue(m.getDiff(ioStat.MergedReadCount, lastMeasure.MergedReadCount)))
	data = append(data, formatIOValue(m.getDiff(ioStat.MergedWriteCount, lastMeasure.MergedWriteCount)))
	data = append(data, formatIOValue(readTime))
	data = append(data, formatIOValue(writeTime))
	data = append(data, formatIOValue(ioTime))
	data = append(data, formatIOValue(readBytes))
	data = append(data, formatIOValue(writeBytes))

	lastIOMeasure[m.Device] = getLastIOMeasure(ioStat)

	// ReadCount|WriteCount|MergedReadCount|MergedWriteCount|ReadTime|WriteTime|IoTime|ReadBytes|WriteBytes
	return data, nil
}

func (m *DiskIOStatProvider) GetAverageRecord(records [][]string) []string {
	return m.getAverageRecord(records, m.GetFieldsCount(), false, m.GetEmptyRecordValue)
}

func (m *DiskIOStatProvider) GetFieldsCount() int {
	return 9
}

func (m *DiskIOStatProvider) CheckRecord(data []string, filter StatProviderFilter) bool {
	if filter == nil {
		return true
	}
	return filter.Check(data)
}

func (m *DiskIOStatProvider) GetCode() string {
	return DISK_IO_PROVIDER_CODE + m.Device
}

func (m *DiskIOStatProvider) getDiff(new, old uint64) uint64 {
	if new < old {
		return 0
	}
	return new - old
}

func getLastIOMeasure(ioStat disk.IOCountersStat) IOMeasure {
	return IOMeasure{
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
}

func formatSpaceValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}

func formatIOValue(value uint64) string {
	return strconv.FormatUint(value, 10)
}
