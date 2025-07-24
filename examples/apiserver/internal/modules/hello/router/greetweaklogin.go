package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/middlewares"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/handler"
)

func SetupGreetWeakLoginRoutes(eng *gin.Engine, groupPath string, h *handler.Instance) {
	g := eng.Group(groupPath, middlewares.WeaAuth())
	{
		g.POST("/greet/sayHi", h.SayHi)
	}
}
