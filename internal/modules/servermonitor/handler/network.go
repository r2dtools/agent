package handler

import (
	"github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

func LoadNetworkTimeLineData(requestData *agentintegration.ServerMonitorStatisticsRequestData, logger logger.LoggerInterface) (*agentintegration.ServerMonitorNetworkResponseData, error) {
	var responseData agentintegration.ServerMonitorNetworkResponseData
	responseData.TimeLineData = make(map[string][]agentintegration.ServerMonitorTimeLinePoint)
	responseData.InterfacesInfo = make([]map[string]string, 0)
	filter := &service.StatProviderTimeFilter{
		FromTime: requestData.FromTime,
		ToTime:   requestData.ToTime,
	}

	if err := loadOverallNetworkTimeLineData(&responseData, filter, logger); err != nil {
		return nil, err
	}
	if err := loadNetworkInterfaceInfo(&responseData); err != nil {
		return nil, err
	}

	return &responseData, nil
}

func loadOverallNetworkTimeLineData(responseData *agentintegration.ServerMonitorNetworkResponseData, filter service.StatProviderFilter, logger logger.LoggerInterface) error {
	overallNetworkStatCollector, err := service.GetStatCollector(&service.OverallNetworkStatProvider{}, logger)
	if err != nil {
		return err
	}

	rows, err := overallNetworkStatCollector.Load(filter)
	if err != nil {
		return err
	}

	var overallNetworkData []agentintegration.ServerMonitorTimeLinePoint
	for _, row := range rows {
		overallNetworkData = append(overallNetworkData, getNetworkTimeLinePoint(row))
	}
	responseData.TimeLineData["overall"] = overallNetworkData

	return nil
}

func loadNetworkInterfaceInfo(responseData *agentintegration.ServerMonitorNetworkResponseData) error {
	interfacesInfo, err := service.GetNetworkInterfacesInfo()
	if err != nil {
		return err
	}

	responseData.InterfacesInfo = interfacesInfo
	return nil
}

func getNetworkTimeLinePoint(row []string) agentintegration.ServerMonitorTimeLinePoint {
	return agentintegration.ServerMonitorTimeLinePoint{
		Time: row[0],
		Value: map[string]string{
			"bytesrecv":   row[1],
			"bytessent":   row[2],
			"packetsrecv": row[3],
			"packetssent": row[4],
		},
	}
}
