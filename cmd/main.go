package main

import (
	"TaskScheduler/internal/api"
	"TaskScheduler/internal/db"
	"TaskScheduler/internal/server"
	"TaskScheduler/internal/logger"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла, %v", err)
	}
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = "7550"
	}
	dbPath := os.Getenv("TODO_DBFILE")
	if len(dbPath) == 0 {
		dbPath = "./database/scheduler.db"
	}
	err = db.Init(dbPath,logger.NewDb())
	if err != nil {
		log.Fatalf("Can`t open database, %v", err)
	}
	defer db.Close()
	server := server.New(logger.NewServer())
	api.Init(server,logger.NewApi())
	server.Run(port)
}
