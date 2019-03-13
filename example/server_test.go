package example

import (
	"github.com/csby/wsf/server/host"
	"github.com/csby/wsf/types"
)

type Server struct {
}

func (s *Server) Run(starter func(server types.Server)) error {
	httpHandler := &HttpHandler{}
	httpHost := host.NewHost(log, &cfg.Server, httpHandler, nil)
	server, err := host.NewServer(log, httpHost, "server-test")
	if err != nil {
		return err
	}
	if !server.Interactive() {
		httpHandler.restart = server.Restart
	}

	if starter != nil {
		go starter(server)
	}

	return server.Run()
}
