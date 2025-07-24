package handler

import (
	"context"

	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello/service"
)

type Instance struct {
	service *service.Instance
}

func New(ctx context.Context) *Instance {
	return &Instance{
		service: service.New(),
	}
}
