package host

import (
	"github.com/csby/wsf/types"
	"github.com/kardianos/service"
)

type program struct {
	types.Base

	server types.Host
}

func (s *program) Start(svc service.Service) error {
	s.LogInfo("service '", svc.String(), "' started")

	go s.run()

	return nil
}

func (s *program) Stop(svc service.Service) error {
	s.LogInfo("service '", svc.String(), "' stopped")

	return nil
}

func (s *program) run() {
	if s.server == nil {
		return
	}

	err := s.server.Run()
	if err != nil {
		s.LogError(err)
	}
}
