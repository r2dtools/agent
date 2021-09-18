package service

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// MemoryStatPrivider retrieves statistics data for memory
type MemoryStatPrivider struct{}

func (sc *MemoryStatPrivider) GetData() ([]string, error) {
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

	// time|total|available|free|used|active|inactive|cached|buffers
	return data, nil
}

func (sc *MemoryStatPrivider) GetCode() string {
	return MEMORY_PROVIDER_CODE
}

func (sc *MemoryStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if filter == nil {
		return true
	}

	return filter.Check(data)
}

func formatMemValue(value uint64) string {
	return strconv.FormatUint(value, 10)
}
