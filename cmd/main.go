package main

import (
	"TaskScheduler/internal/api"
	"TaskScheduler/internal/db"
	"TaskScheduler/internal/logger"
	"TaskScheduler/internal/server"
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
		port = "7540"
	}
	dbPath := os.Getenv("TODO_DBFILE")
	if len(dbPath) == 0 {
		dbPath = "./database/scheduler.db"
	}
	password := os.Getenv("TODO_PASSWORD")
	if len(password) == 0 {
		password = "1234"
	}
	jwtSecret := os.Getenv("TODO_JWT_SECRET")
	if len(jwtSecret) == 0 {
		jwtSecret = "secret"
	}
	err = db.Init(dbPath, logger.NewDb())
	if err != nil {
		log.Fatalf("Can`t open database, %v", err)
	}
	defer db.Close()
	server := server.New(logger.NewServer())
	api.Init(server, logger.NewApi())
	api.InitPassword(password)
	api.InitJwtSecret(jwtSecret)

	server.Run(port)
}
