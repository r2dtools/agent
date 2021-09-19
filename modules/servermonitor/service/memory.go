package service

import "strconv"

type Memory struct{}

func (m *Memory) formatMemValue(value uint64) string {
	return strconv.FormatUint(value/(1024*1024), 10)
}
