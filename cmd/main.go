package main

import (
	"TaskScheduler/internal/api"
	"TaskScheduler/internal/db"
	"TaskScheduler/internal/server"
	"log"
	"os"

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
	server := server.New()
	api.Init(server)
	server.Run(port)
}
