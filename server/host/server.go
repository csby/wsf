package host

import (
	"fmt"
	"github.com/csby/wsf/types"
	"github.com/kardianos/service"
)

func NewServer(log types.Log, host types.Host, serviceName string) (types.Server, error) {
	instance := &server{}
	instance.SetLog(log)
	instance.program.server = host

	svc, err := service.New(&instance.program, &service.Config{Name: serviceName, DisplayName: serviceName})
	if err != nil {
		return nil, err
	}
	instance.service = svc

	return instance, nil
}

type server struct {
	types.Base

	program program
	service service.Service
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

func (s *server) Restart() error {
	if s.Interactive() {
		return fmt.Errorf("not running as service")
	} else {
		return s.service.Restart()
	}
}

func (s *server) Start() error {
	if s.Interactive() {
		return fmt.Errorf("not running as service")
	} else {
		return s.service.Start()
	}
}

func (s *server) Stop() error {
	if s.Interactive() {
		return s.program.server.Close()
	} else {
		return s.service.Stop()
	}
}

func (s *server) Install() error {
	if s.Interactive() {
		return fmt.Errorf("not running as service")
	} else {
		return s.service.Install()
	}
}

func (s *server) Uninstall() error {
	if s.Interactive() {
		return fmt.Errorf("not running as service")
	} else {
		return s.service.Uninstall()
	}
}

func (s *server) Status() (types.ServerStatus, error) {
	if s.Interactive() {
		return types.ServerStatusUnknown, fmt.Errorf("not running as service")
	}

	status, err := s.service.Status()
	if err != nil {
		return types.ServerStatusUnknown, err
	}

	return types.ServerStatus(status), nil
}
