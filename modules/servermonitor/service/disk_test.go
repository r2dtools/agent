package service

import (
	"strings"
	"testing"

	"github.com/shirou/gopsutil/disk"
	"github.com/unknwon/com"
)

func TestGetDiskAverageRecord(t *testing.T) {
	records := [][]string{
		{`{"1":"4","2":"5"}`},
		{`{"1":"10","2":"5"}`},
		{`{"1":"4","2":"6","3":"2"}`},
		{`{"1":"2","2":"4","3":"4"}`},
	}
	expectedAverage := []string{`{"1":"5","2":"5","3":"3"}`}
	provider := &DiskUsageStatProvider{}
	average := provider.GetAverageRecord(records)
	if !com.CompareSliceStr(average, expectedAverage) {
		t.Errorf("invalid average record. Expected %s, got %s", strings.Join(expectedAverage, ","), strings.Join(average, ","))
	}
}

func TestGetDiskDevicesFromPartitions(t *testing.T) {
	partitions := []disk.PartitionStat{
		{
			Device:     "/dev/sdc",
			Mountpoint: "/",
			Fstype:     "ext4",
		},
		{
			Device:     "/loop/sdc",
			Mountpoint: "/folderc",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/sdb",
			Mountpoint: "/folderb",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/sdb1",
			Mountpoint: "/folderb1",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/sda1",
			Mountpoint: "/foldera1",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/sda2",
			Mountpoint: "/foldera2",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/nvme0n1p1",
			Mountpoint: "/foldernvmep1",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/nvme0n1p2",
			Mountpoint: "/foldernvmep2",
			Fstype:     "ext4",
		},
		{
			Device:     "/dev/dm-1",
			Mountpoint: "/folderdm1",
			Fstype:     "ext4",
		},
	}
	expectedDevices := []string{"sdc", "sdb", "sda1", "sda2", "nvme0n1p1", "nvme0n1p2", "dm-1"}
	devices, err := getDiskDevicesFromPartitions(partitions)
	if err != nil {
		t.Fatal(err)
	}
	if !com.CompareSliceStrU(expectedDevices, devices) {
		t.Errorf("invalid device list. Expected %s, got %s", strings.Join(expectedDevices, ","), strings.Join(devices, ","))
	}
}
