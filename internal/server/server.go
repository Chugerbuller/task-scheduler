package server

import (
	"TaskScheduler/internal/logger"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	Router *chi.Mux
	Logger *logger.Logger
}

func New(logger *logger.Logger) *Server {
	return &Server{
		Router: chi.NewRouter(),
		Logger: logger,
	}
}
func (s *Server) Run(port string) {
	s.Logger.Print("Server started on port: %s", port)
	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		s.Logger.Fatal("Server is shuted down %v", err)
	}
}
