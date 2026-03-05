package config

import (
	"fmt"
	"log"
	"os"
	"rextra-backend/db"
	"rextra-backend/internal/middleware"
	auth "rextra-backend/internal/modules/auth"
	jelajahprofesi "rextra-backend/internal/modules/jelajah_profesi"
	kenalidiri "rextra-backend/internal/modules/kenali_diri"
	persona "rextra-backend/internal/modules/persona"
	token "rextra-backend/internal/modules/token"

	"rextra-backend/internal/pkg/cache"
	"rextra-backend/internal/pkg/export"
	myfirebase "rextra-backend/internal/pkg/firebase"

	"github.com/gin-gonic/gin"
)

type RestConfig struct {
	server       *gin.Engine
	cacheService cache.CacheService
}

func NewRest() RestConfig {
	db := db.New()

	// Mode
	mode := os.Getenv("APP_MODE")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.Default()
	server := NewRouter(app)
	firebaseApp := myfirebase.New()
	middleware := middleware.New(db, firebaseApp.MustGetClient())
	cacheService := cache.New()

	var exportService export.ExportService = export.New()

	// Module
	auth.InitModule(server, db, middleware)
	persona.InitModule(server, db, middleware)
	kenalidiri.InitModule(server, db, middleware, cacheService, exportService)
	token.InitModule(server, db, middleware)
	jelajahprofesi.InitModule(server, db, middleware)

	return RestConfig{
		server:       server,
		cacheService: cacheService,
	}
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")
	if port == "" {
		port = "8000"
	}

	serve := fmt.Sprintf("%s:%s", host, port)
	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
	log.Println("server start on port ", serve)
}

func (ap *RestConfig) Close() error {
	if ap.cacheService != nil {
		if closer, ok := ap.cacheService.(interface{ Close() error }); ok {
			return closer.Close()
		}
	}
	return nil
}
