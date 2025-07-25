package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/internal/handler"
	"github.com/zmjung/jamesdb/internal/router"
)

func main() {
	var cfg config.Config
	err := config.LoadConfig(&cfg, "config.yml")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	engine := gin.Default()

	graphHandler := handler.NewGraphHandler(cfg)
	router := router.NewRouter(graphHandler)
	router.SetupRoutes(engine)

	engine.Run(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port))
}
