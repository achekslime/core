package rest_api

import (
	"github.com/achekslime/core/app"
	"github.com/achekslime/core/gin_server"
	"github.com/gin-gonic/gin"
)

type RestApiRunner struct {
	ginRouter *gin.Engine
	addr      string
}

func NewService() *RestApiRunner {
	return &RestApiRunner{}
}

func (api RestApiRunner) ConfigureServer(router *gin.Engine, addr string) {
	api.ginRouter = router
	api.addr = addr
}

func (api RestApiRunner) Run() {
	serviceWorker := app.NewServerScript()
	serviceWorker.Tasks(gin_server.StartGin(api.ginRouter, api.addr))

	app.StartServer(serviceWorker)
}
