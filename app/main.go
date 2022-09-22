package main

import (
	"os"

	impulse_api "impulse-api"
	"impulse-api/pkg/handler"
	"impulse-api/pkg/repository"
	"impulse-api/pkg/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title Sendpulse Western Horoscope
// @version 1.0
// @description API takes birth_date, birth_time, birth_place and gender. As a response api returns processed data

// @host khusainnov.ru
// @BasePath /signs/{birthday}/{birth_time}/{city}/{sex}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Infof("Config initialization")
	if err := godotenv.Load("./config/.env"); err != nil {
		logrus.Fatalf("Cannot load .env config, due to error: %s", err.Error())
	}

	repo := repository.NewRepository()
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	s := new(impulse_api.Server)

	logrus.Infof("Server starting on port:%s", os.Getenv("PORT"))
	if err := s.Run(os.Getenv("PORT"), handlers.InitRoute()); err != nil {
		logrus.Fatalf("Error due start server: %s", err.Error())
	}
}
