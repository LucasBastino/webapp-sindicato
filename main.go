package main

import (
	"fmt"
	"os"

	"github.com/LucasBastino/webapp-sindicato/internal/bootstrap"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/joho/godotenv"
)


func main() {
	// Loading .env file
	err := godotenv.Load("./.env")
	if err!=nil{
		fmt.Println("failed to load .env: %w", err)
		os.Exit(1)
	}
	
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Println("CRITICAL: invalid config:", err)
		os.Exit(1)
	}

	logger, err := logger.NewSlogLogger()
	if err!=nil{
		fmt.Println("CRITICAL: failed to init logger", "err", err)
		os.Exit(1)
	}
	defer logger.File.Close()

	app, cleanup, err := bootstrap.InitApp(cfg, logger)
	if err!=nil{
		logger.Error("CRITICAL: failed to init app", "err", err)
		os.Exit(1)
	}
	// para cerrar la db
	defer cleanup()
	
	err = app.Listen(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err!=nil{
		logger.Error("server crashed", "err", err)
		os.Exit(1)
	}
	
}
