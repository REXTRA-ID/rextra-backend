package main

import (
	"os"
	"rextra-backend/cmd"
	"rextra-backend/db"
	"rextra-backend/db/migrations"
	seeders "rextra-backend/db/seeder"
	"rextra-backend/internal/config"
	mylog "rextra-backend/internal/pkg/logger"

	"github.com/joho/godotenv"
)

const ENV_FILE = ".env.dev"

func main() {
	_ = godotenv.Load(ENV_FILE)

	// Jika dipanggil dengan flag (CLI mode)
	if len(os.Args) > 1 {
		if err := cmd.Commands(); err != nil {
			panic("Failed Get Commands: " + err.Error())
		}
		return
	}

	// Mode API Server: Jalankan auto-migrate & seeder di Dev Mode
	if os.Getenv("APP_MODE") == "development" {
		database := db.New()
		mylog.Infof("[AUTO] Running auto migration...")
		_ = migrations.Migrate(database)
		mylog.Infof("[AUTO] Running auto seeder...")
		_ = seeders.Seeding(database)
	}

	RestApi := config.NewRest()
	RestApi.Start()
}
