package host

import (
	"fmt"
	"github.com/csby/wsf/types"
	"github.com/kardianos/service"
)

func NewServer(log types.Log, host types.Host, serviceName string, serviceArguments ...string) (types.Server, error) {
	instance := &server{}
	instance.SetLog(log)
	instance.program.server = host
	instance.serviceName = serviceName

	cfg := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
	}
	if len(serviceArguments) > 0 {
		cfg.Arguments = serviceArguments
	}
	svc, err := service.New(&instance.program, cfg)
	if err != nil {
		return nil, err
	}
	instance.service = svc

	return instance, nil
}

type server struct {
	types.Base

	program     program
	service     service.Service
	serviceName string
}

func (s *server) ServiceName() string {
	return s.serviceName
}

func (s *server) Interactive() bool {
	return service.Interactive()
}

func (s *server) Run() error {
	if s.program.server == nil {
		return fmt.Errorf("invalid host: nil")
	}

	if s.Interactive() {
		return s.program.server.Run()
	} else {
		return s.service.Run()
	}
}

func (s *server) Shutdown() error {
	if s.program.server == nil {
		return fmt.Errorf("invalid host: nil")
	}

	return s.program.server.Close()
}

func (s *server) Restart() error {
	return s.service.Restart()
}

func (s *server) Start() error {
	return s.service.Start()
}

func (s *server) Stop() error {
	return s.service.Stop()
}

func (s *server) Install() error {
	return s.service.Install()
}

func (s *server) Uninstall() error {
	err := s.service.Stop()
	if err != nil {
	}

	return s.service.Uninstall()
}

func (s *server) Status() (types.ServerStatus, error) {
	status, err := s.service.Status()
	if err != nil {
		return types.ServerStatusUnknown, err
	}

	return types.ServerStatus(status), nil
}
