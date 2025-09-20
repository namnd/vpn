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
	r.PUT("/start-node/:id", routes.StartNode)
	r.PUT("/stop-node/:id", routes.StopNode)
	r.DELETE("/delete-node/:id", routes.DeleteNode)

	r.Run(":8080")
}
