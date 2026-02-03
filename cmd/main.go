package main

import (
	"TaskScheduler/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла, %v", err)
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7550"
	}
	dbPath := os.Getenv("TODO_dbPath")
	storage, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Can`t open database, %v", err)
	}
	defer storage.Close()

	router := chi.NewRouter()

	router.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Printf("Server started on port:%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server is shuted down %v", err)
	}
}
