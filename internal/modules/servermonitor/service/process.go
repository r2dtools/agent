package service

import (
	"github.com/r2dtools/agentintegration"
	proc "github.com/shirou/gopsutil/process"
)

func GetProcessesData() ([]agentintegration.ProcessData, error) {
	processes, err := proc.Processes()
	if err != nil {
		return nil, err
	}

	var data []agentintegration.ProcessData
	for _, process := range processes {
		var processData agentintegration.ProcessData

		name, err := process.Name()
		if err != nil {
			continue
		}
		processData.Name = name

		processData.Pid = process.Pid
		ppid, err := process.Ppid()
		if err != nil {
			continue
		}
		processData.PPid = ppid

		user, err := process.Username()
		if err != nil {
			continue
		}
		processData.User = user
		cmd, err := process.Cmdline()
		if err != nil {
			continue
		}
		processData.Cmd = cmd
		cpu, err := process.CPUPercent()
		if err != nil {
			continue
		}
		processData.Cpu = cpu

		memory, err := process.MemoryPercent()
		if err != nil {
			continue
		}
		processData.Memory = memory
		openFiles, err := process.OpenFiles()
		if err != nil {
			continue
		}
		for _, openFile := range openFiles {
			processData.OpenFiles = append(processData.OpenFiles, openFile.Path)
		}

		netStats, err := process.NetIOCounters(false)
		if err != nil || len(netStats) == 0 {
			continue
		}
		netStat := netStats[0]
		processData.NetBytesRecv = netStat.BytesRecv
		processData.NetBytesSent = netStat.BytesSent
		processData.NetPacketsRecv = netStat.PacketsRecv
		processData.NetPacketsSent = netStat.PacketsSent

		ioStat, err := process.IOCounters()
		if err != nil {
			continue
		}
		processData.DiskReadBytes = ioStat.ReadBytes
		processData.DiskWriteBytes = ioStat.WriteBytes
		processData.DiskReadCount = ioStat.ReadCount
		processData.DiskWriteCount = ioStat.WriteCount

		data = append(data, processData)
	}

	return data, nil
}
