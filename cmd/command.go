package cmd

import (
	"fmt"
	"rextra-backend/db"
	"rextra-backend/db/migrations"
	seeders "rextra-backend/db/seeder"
	mylog "rextra-backend/internal/pkg/logger"

	"os"
	"os/exec"
	"runtime"

	"gorm.io/gorm"
)

func Commands() error {
	db := db.New()
	if err := getParams(db); err != nil {
		return err
	}

	return nil
}

func getParams(db *gorm.DB) error {
	migrate := false
	seeder := false
	watch := false
	reset := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}
		if arg == "--seeder" {
			seeder = true
		}
		if arg == "--watch" {
			watch = true
		}
		if arg == "--reset-db" {
			reset = true
		}
	}
	if migrate {
		if err := migrations.Migrate(db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		mylog.Infof("Migration completed successfully")
	}

	if seeder {
		if err := seeders.Seeding(db); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
		mylog.Infof("Seeder completed successfully")
	}

	if watch {
		if err := runWatch(); err != nil {
			return fmt.Errorf("watching failed: %w", err)
		}
		mylog.Infof("Start watching program")
		os.Exit(0)
	}

	if reset {
		if err := resetDatabase(db); err != nil {
			return fmt.Errorf("reseting failed: %w", err)
		}
	}

	if seeder || watch || migrate || reset {
		os.Exit(0)
	}

	return nil
}

func runWatch() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "air -c .air.windows.toml")
	case "linux", "darwin":
		cmd = exec.Command("bash", "-c", "air -c .air.linux.toml")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		mylog.Errorf("Error running command: %s", err)
		return err
	}

	mylog.Infoln("Command executed successfully")
	return nil
}

func resetDatabase(db *gorm.DB) error {

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Reseting ==========="))
	mylog.Infof("Reseting Tables...")

	tables, err := db.Migrator().GetTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		sql := fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE`, table)

		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}

	mylog.Infof("Reset completed...")

	return nil
}
