package hello

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/handler"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/router"
	"github.com/mel0dys0ng/song/pkgs/https"
)

func SetupRoutes(groupPath string) https.Option {
	handler := handler.New(context.Background())
	return https.Routes(
		setupGreetRoutes(groupPath, handler),
	)
}

func setupGreetRoutes(groupPath string, handler *handler.Instance) https.Route {
	return func(eng *gin.Engine) {
		router.SetupGreetMustLoginRoutes(eng, groupPath, handler)
		router.SetupGreetWeakLoginRoutes(eng, groupPath, handler)
	}
}
