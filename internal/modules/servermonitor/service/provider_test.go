package service

import (
	"strings"
	"testing"

	"github.com/unknwon/com"
)

func TestGetAverageRecord(t *testing.T) {
	records := [][]string{
		{"1", "2", "0", "10", "15"},
		{"3", "2", "1", "16", ""},
		{"6", "2", "5", "12", ""},
		{"2", "2", "2", "22", "15"},
	}
	expectedAverageInt := []string{"3", "2", "2", "15", "15"}
	expectedAverageFloat := []string{"3.00", "2.00", "2.00", "15.00", "15.00"}
	provider := &OverallCPUStatPrivider{}

	averageInt := provider.getAverageRecord(records, 5, false, func(index int) string { return "" })
	if !com.CompareSliceStr(averageInt, expectedAverageInt) {
		t.Errorf("invalid average record. Expected %s, got %s", strings.Join(expectedAverageInt, ","), strings.Join(averageInt, ","))
	}

	averageFloat := provider.getAverageRecord(records, 5, true, func(index int) string { return "" })
	if !com.CompareSliceStr(averageFloat, expectedAverageFloat) {
		t.Errorf("invalid average record. Expected %s, got %s", strings.Join(expectedAverageFloat, ","), strings.Join(averageFloat, ","))
	}
}
