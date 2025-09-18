package main

import (
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/routes"
)

func main() {

	r := gin.Default()

	r.GET("/", routes.Home)

	r.Run(":8080")
}
