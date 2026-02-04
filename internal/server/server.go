package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	Router *chi.Mux
}

func New() *Server {
	return &Server{
		Router: chi.NewRouter(),
	}
}
func (s *Server) Run(port string) {
	log.Printf("Server started on port:%s", port)

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		log.Fatalf("Server is shuted down %v", err)
	}
}
