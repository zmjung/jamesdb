package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/internal/router"
	"strconv"
)

func main() {
	var cfg config.Config
	err := config.LoadConfig(&cfg, "config.yml")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	
	engine := gin.Default()
	router.SetupRoutes(engine)
	engine.Run(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port))
}
