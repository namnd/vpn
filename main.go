package main

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/routes"
)

//go:embed assets/*
var f embed.FS

func main() {

	r := gin.Default()
	r.StaticFS("/public", http.FS(f))

	r.GET("/", routes.Home)
	r.GET("/:country", routes.Home)

	r.POST("/add-new-node", routes.CreateNode)
	r.POST("/start-node/:id", routes.StartNode)
	r.POST("/stop-node/:id", routes.StopNode)

	r.Run(":8080")
}
