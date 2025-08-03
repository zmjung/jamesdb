package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/internal/grapher"
	"github.com/zmjung/jamesdb/internal/handler"
	"github.com/zmjung/jamesdb/internal/middleware"
	"github.com/zmjung/jamesdb/internal/router"
)

func main() {
	var cfg config.Config
	err := config.LoadConfig(&cfg, "config.yml")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	engine := gin.Default()
	engine.Use(middleware.GetLogging())
	engine.Use(middleware.GetRecovery())

	gs := grapher.NewGraphService(&cfg)
	if gs == nil {
		panic("Failed to create GraphWriter")
	}
	graphHandler := handler.NewGraphHandler(cfg, gs)
	router := router.NewRouter(graphHandler)
	router.SetupRoutes(engine)

	engine.Run(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port))
}
