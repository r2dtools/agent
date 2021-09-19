package service

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// VirtualMemoryStatPrivider retrieves statistics data for memory
type VirtualMemoryStatPrivider struct{}

func (m *VirtualMemoryStatPrivider) GetData() ([]string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatMemValue(vmStat.Total))
	data = append(data, formatMemValue(vmStat.Available))
	data = append(data, formatMemValue(vmStat.Free))
	data = append(data, formatMemValue(vmStat.Used))
	data = append(data, formatMemValue(vmStat.Active))
	data = append(data, formatMemValue(vmStat.Inactive))
	data = append(data, formatMemValue(vmStat.Cached))
	data = append(data, formatMemValue(vmStat.Buffers))

	// time|total|available|free|used|active|inactive|cached|buffered
	return data, nil
}

func (m *VirtualMemoryStatPrivider) GetCode() string {
	return VIRTUAL_MEMORY_PROVIDER_CODE
}

func (m *VirtualMemoryStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 9 {
		return false
	}

	if filter == nil {
		return true
	}

	return filter.Check(data)
}

// SwapMemoryStatPrivider retrieves statistics data for memory
type SwapMemoryStatPrivider struct{}

func (m *SwapMemoryStatPrivider) GetData() ([]string, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, formatMemValue(swapStat.Total))
	data = append(data, formatMemValue(swapStat.Used))
	data = append(data, formatMemValue(swapStat.Free))

	// time|total|used|free
	return data, nil
}

func (m *SwapMemoryStatPrivider) GetCode() string {
	return SWAP_MEMORY_PROVIDER_CODE
}

func (m *SwapMemoryStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 4 {
		return false
	}

	if filter == nil {
		return true
	}

	return filter.Check(data)
}

func formatMemValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}
