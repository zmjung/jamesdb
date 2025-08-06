package main

import (
	"flag"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/internal/grapher"
	"github.com/zmjung/jamesdb/internal/handler"
	"github.com/zmjung/jamesdb/internal/log"
	"github.com/zmjung/jamesdb/internal/middleware"
	"github.com/zmjung/jamesdb/internal/router"
)

func loadConfig() *config.Config {
	configPath := flag.String("config", "config.yml", "configuration file path")
	flag.Parse()

	cfg := &config.Config{}
	err := config.LoadConfig(cfg, *configPath)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	return cfg
}

func main() {
	cfg := loadConfig()
	log.SetDefaultLogger(cfg)
	slog.Debug("Using config file", "config", cfg)

	engine := gin.Default()
	engine.Use(middleware.GetLogging())
	engine.Use(middleware.GetRecovery())

	gs := grapher.NewGraphService(cfg)
	if gs == nil {
		panic("Failed to create GraphWriter")
	}
	graphHandler := handler.NewGraphHandler(cfg, gs)
	router := router.NewRouter(graphHandler)
	router.SetupRoutes(engine)

	if engine.Run(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port)) != nil {
		panic("Failed to start gin engine")
	}
}
