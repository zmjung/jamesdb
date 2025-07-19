package main

import (
	"github.com/zmjung/jamesdb/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	router.SetupRoutes(engine)
	engine.Run(":8080")
}
