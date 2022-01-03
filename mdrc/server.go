package mdrc

import (
	"net"
	"net/http"

	"github.com/go-logr/logr"
)

type Server struct {
	logger     logr.Logger
	controller *Controller
}

func NewServer(l logr.Logger, c *Controller) *Server {
	return &Server{
		logger:     l.WithName("server"),
		controller: c,
	}
}

func (s *Server) Serve(port string) {
	const network = "tcp"
	s.logger.Info("listening...", "network", network, "port", port)
	l, err := net.Listen(network, ":"+port)
	if err != nil {
		s.logger.Error(err, "unable to listen")
		return
	}

	s.logger.Info("serving...", "url", "http://localhost:"+port)
	http.HandleFunc("/", s.controller.HandleHTML())
	http.HandleFunc("/run/", s.controller.HandleCommand())
	if err := http.Serve(l, nil); err != nil {
		s.logger.Error(err, "unable to serve")
		return
	}
}
