package main

import (
	"rextra-backend/cmd"
	"rextra-backend/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := cmd.Commands(); err != nil {
		panic("Failed Get Commands: " + err.Error())
	}

	RestApi := config.NewRest()
	defer RestApi.Close()
	RestApi.Start()
}
