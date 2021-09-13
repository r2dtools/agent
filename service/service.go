package service

import (
	"fmt"
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
func (s *ServiceManager) RunServices() error {
	var errs []error

	for name, service := range s.services {
		runService := func(iName string, iService Service, errs []error) {
			if err := iService.Run(); err != nil {
				_ = append(errs, fmt.Errorf("could not start service '%s': %v", iName, err))
			}
		}

		go runService(name, service, errs)
	}

	if len(errs) != 0 {
		return errs[0]
	}

	return nil
}
