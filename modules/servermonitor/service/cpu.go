package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type OverallCPUStatPrivider struct{}

func (sc *OverallCPUStatPrivider) GetData() ([]string, error) {
	items, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	var data []string
	item := items[0]
	total := item.Total()
	data = append(data, strconv.FormatInt(time.Now().Unix(), 10))
	data = append(data, sc.getPercertValue(item.System, total))
	data = append(data, sc.getPercertValue(item.User, total))
	data = append(data, sc.getPercertValue(item.Nice, total))
	data = append(data, sc.getPercertValue(item.Idle, total))

	// data: time|system|user|nice|idle
	return data, nil
}

func (sc *OverallCPUStatPrivider) GetCode() string {
	return OVERALL_CPU_PROVIDER_CODE
}

func (sc *OverallCPUStatPrivider) getPercertValue(value, total float64) string {
	return fmt.Sprintf("%.2f", 100*value/total)
}

func (sc *OverallCPUStatPrivider) CheckData(data []string, filter StatProviderFilter) bool {
	if len(data) != 5 {
		return false
	}

	if filter != nil {
		return filter.Check(data)
	}

	return true
}
