package service

import (
	"strings"
	"testing"
	"time"

	"github.com/unknwon/com"
)

func TestGetAverageCount(t *testing.T) {
	items := [][]int{
		{1636256509, 1636285309, -1},
		{1635783176, 1636301576, 5},
		{1635610242, 1636301442, 10},
	}
	statCollector := &StatCollector{}
	for _, item := range items {
		average := statCollector.getAverageCount(item[0], item[1])
		if item[2] != average {
			t.Errorf("invalid average count. Expected %d, got %d", item[2], average)
		}
	}
}

func TestGetFullRecord(t *testing.T) {
	type itemData struct {
		record,
		expectedRecord []string
		fieldsCount int
	}
	statCollector := StatCollector{Provider: &OverallCPUStatPrivider{}}
	items := []itemData{
		{record: []string{"0.39", "0.35", "0", "99.18", "0.82"}, expectedRecord: []string{"0.39", "0.35", "0", "99.18", "0.82"}, fieldsCount: 5},
		{record: []string{"0.39", "0.35", "0"}, expectedRecord: []string{"0.39", "0.35", "0", "", ""}, fieldsCount: 5},
		{record: []string{"0.39", "0.35", "0", "99.18", "0.82"}, expectedRecord: []string{"0.39", "0.35", "0"}, fieldsCount: 3},
		{record: []string{"0.39"}, expectedRecord: []string{"0.39", "", ""}, fieldsCount: 3},
	}

	for _, item := range items {
		record := statCollector.getFullRecord(item.fieldsCount, item.record)
		if !com.CompareSliceStr(record, item.expectedRecord) {
			t.Errorf("invalid full record. Expected %s, got %s", strings.Join(item.expectedRecord, ","), strings.Join(record, ","))
		}
	}
}

func TestGetExtendedRecords(t *testing.T) {
	statCollector := StatCollector{Provider: &OverallCPUStatPrivider{}}
	type itemData struct {
		record,
		previousRecord []string
		interval        time.Duration
		expectedRecords [][]string
	}
	items := []itemData{
		{
			record:         []string{"1635073416", "0.39", "0.35", "0", "99.18", "0.82"},
			previousRecord: nil,
			interval:       time.Minute,
			expectedRecords: [][]string{
				{"1635073416", "0.39", "0.35", "0", "99.18", "0.82"},
			},
		},
		{
			record:         []string{"1635073416", "0.39", "0.35", "0", "99.18", "0.82"},
			previousRecord: []string{"1635073366", "0.39", "0.35", "0", "99.18", "0.82"},
			interval:       time.Minute,
			expectedRecords: [][]string{
				{"1635073416", "0.39", "0.35", "0", "99.18", "0.82"},
			},
		},
		{
			record:         []string{"1635073416", "0.39", "0.35", "0"},
			previousRecord: []string{"1635073200", "0.39", "0.35", "0", "99.18", "0.82"},
			interval:       time.Minute,
			expectedRecords: [][]string{
				{"1635073260", "", "", "", "", ""},
				{"1635073320", "", "", "", "", ""},
				{"1635073380", "", "", "", "", ""},
				{"1635073416", "0.39", "0.35", "0", "", ""},
			},
		},
	}
	for _, item := range items {
		extendedRecords := statCollector.getExtendedRecords(item.previousRecord, item.record, item.interval)
		for index, extendedRecord := range extendedRecords {
			if !com.CompareSliceStr(extendedRecord, item.expectedRecords[index]) {
				t.Errorf("invalid extended record. Expected %s, got %s", strings.Join(item.expectedRecords[index], ","), strings.Join(extendedRecord, ","))
			}
		}
	}
}
