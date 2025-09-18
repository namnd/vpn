package routes

import (
	"net/http"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/ui"
)

func Home(c *gin.Context) {

	p := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		ui.Home(),
	)

	c.Render(http.StatusOK, p)
}
