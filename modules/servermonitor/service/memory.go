package service

import (
	"strconv"

	"github.com/shirou/gopsutil/mem"
)

// VirtualMemoryStatPrivider retrieves statistics data for memory
type VirtualMemoryStatPrivider struct {
	BaseStatProvider
}

func (m *VirtualMemoryStatPrivider) GetRecord() ([]string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, formatMemValue(vmStat.Total))
	data = append(data, formatMemValue(vmStat.Available))
	data = append(data, formatMemValue(vmStat.Free))
	data = append(data, formatMemValue(vmStat.Used))
	data = append(data, formatMemValue(vmStat.Active))
	data = append(data, formatMemValue(vmStat.Inactive))
	data = append(data, formatMemValue(vmStat.Cached))
	data = append(data, formatMemValue(vmStat.Buffers))

	// total|available|free|used|active|inactive|cached|buffered
	return data, nil
}

func (m *VirtualMemoryStatPrivider) GetAverageRecord(records [][]string) []string {
	return m.getAverageRecord(records, m.GetFieldsCount(), false, m.GetEmptyRecordValue)
}

func (m *VirtualMemoryStatPrivider) GetFieldsCount() int {
	return 8
}

func (m *VirtualMemoryStatPrivider) GetCode() string {
	return VIRTUAL_MEMORY_PROVIDER_CODE
}

func (m *VirtualMemoryStatPrivider) CheckRecord(data []string, filter StatProviderFilter) bool {
	if filter == nil {
		return true
	}

	return filter.Check(data)
}

// SwapMemoryStatPrivider retrieves statistics data for memory
type SwapMemoryStatPrivider struct {
	BaseStatProvider
}

func (m *SwapMemoryStatPrivider) GetRecord() ([]string, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	// if there is no SWAP skip statistics collection
	if swapStat.Total == 0 {
		return nil, nil
	}

	var data []string
	data = append(data, formatMemValue(swapStat.Total))
	data = append(data, formatMemValue(swapStat.Used))
	data = append(data, formatMemValue(swapStat.Free))

	// total|used|free
	return data, nil
}

func (m *SwapMemoryStatPrivider) GetAverageRecord(records [][]string) []string {
	return m.getAverageRecord(records, m.GetFieldsCount(), false, m.GetEmptyRecordValue)
}

func (m *SwapMemoryStatPrivider) GetFieldsCount() int {
	return 3
}

func (m *SwapMemoryStatPrivider) GetCode() string {
	return SWAP_MEMORY_PROVIDER_CODE
}

func (m *SwapMemoryStatPrivider) CheckRecord(data []string, filter StatProviderFilter) bool {
	if filter == nil {
		return true
	}

	return filter.Check(data)
}

func formatMemValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}
