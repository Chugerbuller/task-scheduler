package api

import (
	"TaskScheduler/internal/server"
	"net/http"
)

func HandleIndex(w *http.ResponseWriter, r *http.Request) {

}
func Init(s *server.Server) {
	s.Router.Get("/*", http.FileServer(http.Dir("./web")).ServeHTTP)
	s.Router.Get("/api/nextdate", nextDayHandler)

}
