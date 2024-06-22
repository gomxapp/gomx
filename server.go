package gomx

import (
	"errors"
	"fmt"
	"github.com/gomxapp/gomx/internal/config"
	"net/http"
)

type Server struct {
	s *http.Server
	r *Router
}

func init() {
	config.Init()
}

// NewServer initializes a new Server instance and sets server.Server.handler = server.Router.
func NewServer(s *http.Server, r *Router) *Server {
	s.Handler = r
	newServer := &Server{
		s: s,
		r: r,
	}
	return newServer
}

// ListenAndServe wraps http.Server.ListenAndServe. It checks if the given
// router has been initialized.
func (server *Server) ListenAndServe() error {
	defer func() {
		if r := recover(); r != nil {
			server.r.Init()
			_ = server.ListenAndServe()
		}
	}()
	if !server.r.IsInitialized() {
		return errors.New("router was not initialized")
	}
	fmt.Println("Listening at " + server.s.Addr)
	err := server.s.ListenAndServe()
	return err
}

func (server *Server) Close() error {
	return server.s.Close()
}
