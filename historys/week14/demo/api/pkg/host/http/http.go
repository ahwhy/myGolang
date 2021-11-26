package http

import (
	"fmt"

	"github.com/julienschmidt/httprouter"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/ahwhy/myGolang/week14/demo/api/pkg"
	"github.com/ahwhy/myGolang/week14/demo/api/pkg/host"
)

type handler struct {
	service host.Service
	log     logger.Logger
}

var (
	api = &handler{}
)

func (h *handler) Config() error {
	h.log = zap.L().Named("Host")
	if pkg.Host == nil {
		return fmt.Errorf("dependence service host not ready")
	}

	h.service = pkg.Host

	return nil
}

func RegistAPI(r *httprouter.Router) {
	api.Config()

	r.POST("/hosts", api.CreateHost)

	r.PATCH("/hosts/:id", api.PatchHost)
	r.PUT("/hosts/:id", api.PutHost)

	r.GET("/hosts", api.QueryHost)
	r.GET("/hosts/:id", api.DescribeHost)

	r.DELETE("/hosts/:id", api.DeleteHost)
}