package service

import (
	"fmt"

	"github.com/r2dtools/agent/logger"
)

// Service is implemeted by services that are started with starting agent tcp server
type Service interface {
	Run() error
}

type ServiceManager struct {
	services map[string]Service
}

// AddService adds new service to the manager
func (s *ServiceManager) AddService(name string, service Service) {
	if s.services == nil {
		s.services = make(map[string]Service)
	}
	s.services[name] = service
}

// RunServices runs all registered services
func (s *ServiceManager) RunServices() {
	for name, service := range s.services {
		runService := func(iName string, iService Service) {
			if err := iService.Run(); err != nil {
				logger.Error(fmt.Sprintf("could not run service '%s': %v", iName, err))
			}
		}

		go runService(name, service)
	}
}
