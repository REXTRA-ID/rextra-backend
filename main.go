package main

import (
	"rextra-backend/cmd"
	"rextra-backend/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to loading env file")
	}

	if err := cmd.Commands(); err != nil {
		panic("Failed Get Commands: " + err.Error())
	}

	RestApi := config.NewRest()
	RestApi.Start()
}
