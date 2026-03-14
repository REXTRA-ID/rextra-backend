package main

import (
	"rextra-backend/cmd"
	"rextra-backend/internal/config"

	"github.com/joho/godotenv"
)

const ENV_FILE = ".env.dev"

func main() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		panic("Failed to loading env file")
	}

	if err := cmd.Commands(); err != nil {
		panic("Failed Get Commands: " + err.Error())
	}

	RestApi := config.NewRest()
	RestApi.Start()
}
