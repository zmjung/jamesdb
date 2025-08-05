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

func main() {
	configPath := flag.String("config", "config.yml", "configuration file path")
	flag.Parse()

	var cfg config.Config
	err := config.LoadConfig(&cfg, *configPath)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	log.SetDefaultLogger(cfg)
	slog.Debug("Using config file", "configPath", *configPath, "config", cfg)

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

	engine.Run(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port))
}
