package service

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// SwapMemoryStatPrivider retrieves statistics data for memory
type SwapMemoryStatPrivider struct {
	Memory
}

func (m *SwapMemoryStatPrivider) GetData() ([]string, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	var data []string
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, m.formatMemValue(swapStat.Total))
	data = append(data, m.formatMemValue(swapStat.Used))
	data = append(data, m.formatMemValue(swapStat.Free))

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
