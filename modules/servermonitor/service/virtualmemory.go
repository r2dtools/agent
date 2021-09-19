package service

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// VirtualMemoryStatPrivider retrieves statistics data for memory
type VirtualMemoryStatPrivider struct {
	Memory
}

func (m *VirtualMemoryStatPrivider) GetData() ([]string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, m.formatMemValue(vmStat.Total))
	data = append(data, m.formatMemValue(vmStat.Available))
	data = append(data, m.formatMemValue(vmStat.Free))
	data = append(data, m.formatMemValue(vmStat.Used))
	data = append(data, m.formatMemValue(vmStat.Active))
	data = append(data, m.formatMemValue(vmStat.Inactive))
	data = append(data, m.formatMemValue(vmStat.Cached))
	data = append(data, m.formatMemValue(vmStat.Buffers))

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
