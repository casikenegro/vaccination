package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"vaccination-server/db"
	"vaccination-server/models"
	"vaccination-server/router"
	"vaccination-server/server"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}
	db.DBConnection(&db.ConfigDB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
	})
	db.DB.AutoMigrate(models.User{}, models.Drug{}, models.Vaccination{})
	s, err := server.NewServer(context.Background(), &server.Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	})

	if err != nil {
		log.Fatalf("Error creating server %v\n", err)
	}

	s.Start(router.BindRoutes)
}
