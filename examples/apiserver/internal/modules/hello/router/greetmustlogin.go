package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/middlewares"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/handler"
)

func SetupGreetMustLoginRoutes(eng *gin.Engine, groupPath string, h *handler.Instance) {
	g := eng.Group(groupPath, middlewares.MustAuth())
	{
		g.POST("/greet/sayHello", h.SayHello)
	}
}
