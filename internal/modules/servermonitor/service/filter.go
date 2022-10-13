package service

import "strconv"

type StatProviderFilter interface {
	GetFromTime() int
	GetToTime() int
	Check(row []string) bool
}

type StatProviderTimeFilter struct {
	FromTime, ToTime int
}

func (f *StatProviderTimeFilter) Check(row []string) bool {
	if len(row) == 0 {
		return false
	}

	time, err := strconv.Atoi(row[0])
	if err != nil {
		return false
	}

	if f.FromTime > 0 && time < f.FromTime {
		return false
	}

	if f.ToTime > 0 && time > f.ToTime {
		return false
	}

	return true
}

func (f *StatProviderTimeFilter) GetFromTime() int {
	return f.FromTime
}

func (f *StatProviderTimeFilter) GetToTime() int {
	return f.ToTime
}
