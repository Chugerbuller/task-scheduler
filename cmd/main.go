package main

import (
	"TaskScheduler/internal/api"
	"TaskScheduler/internal/db"
	"TaskScheduler/internal/logger"
	"TaskScheduler/internal/server"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен: в Docker и CI переменные приходят из окружения.
	if err := godotenv.Load("./.env"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("Ошибка загрузки .env файла, %v", err)
	}

	port := getEnv("TODO_PORT", "7540")
	dbPath := getEnv("TODO_DBFILE", "./database/scheduler.db")

	// Пустой пароль отключает авторизацию (так же решает middleware auth).
	password := os.Getenv("TODO_PASSWORD")
	jwtSecret := os.Getenv("TODO_JWT_SECRET")
	if password != "" && jwtSecret == "" {
		jwtSecret = randomSecret()
		log.Print("TODO_JWT_SECRET не задан: сгенерирован случайный секрет, токены станут недействительны после перезапуска")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("Can`t create database directory, %v", err)
	}
	if err := db.Init(dbPath, logger.NewDb()); err != nil {
		log.Fatalf("Can`t open database, %v", err)
	}
	defer db.Close()

	server := server.New(logger.NewServer())
	api.Init(server, logger.NewApi())
	api.InitPassword(password)
	api.InitJwtSecret(jwtSecret)

	server.Run(port)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func randomSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Can`t generate JWT secret, %v", err)
	}
	return hex.EncodeToString(b)
}
