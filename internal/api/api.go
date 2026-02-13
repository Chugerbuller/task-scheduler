package api

import (
	"TaskScheduler/internal/logger"
	"TaskScheduler/internal/server"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

var logs *logger.Logger

func Init(s *server.Server, logger *logger.Logger) {
	logs = logger
	s.Router.Post("/api/signin", logging(signIn))
	s.Router.Get("/*", http.FileServer(http.Dir("./web")).ServeHTTP)
	s.Router.Get("/api/nextdate", logging(nextDay))
	s.Router.Post("/api/task", logging(auth(addTask)))
	s.Router.Post("/api/task/done", logging(auth(taskDone)))
	s.Router.Get("/api/tasks", logging(auth(tasks)))
	s.Router.Get("/api/task", logging(auth(task)))
	s.Router.Put("/api/task", logging(auth(updateTask)))
	s.Router.Delete("/api/task", logging(auth(deleteTask)))
}
func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	writeJson(w, errorResponse{
		Error: err.Error(),
	})
}
func writeJson(w http.ResponseWriter, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		logs.PrintErr(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, err := w.Write([]byte(response)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func jsonDeserialize(stream io.ReadCloser, dest any) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return err
	}
	defer stream.Close()
	if err := json.Unmarshal(buf.Bytes(), &dest); err != nil {
		return err
	}
	logs.Print("Password: %s", password)
	return nil
}
