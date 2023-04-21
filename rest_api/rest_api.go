package rest_api

import (
	"fmt"
	"github.com/achekslime/core/app"
	"github.com/achekslime/core/gin_server"
	"github.com/gin-gonic/gin"
)

type RestApiRunner struct {
	ginRouter *gin.Engine
	port      string
}

func NewService() *RestApiRunner {
	return &RestApiRunner{}
}

func (api RestApiRunner) ConfigureServer(router *gin.Engine, port string) {
	api.ginRouter = router
	api.port = port
}

func (api RestApiRunner) Run() {
	serviceWorker := app.NewServerScript()
	serviceWorker.Tasks(gin_server.StartGin(api.ginRouter, fmt.Sprintf(":%s", api.port)))

	app.StartServer(serviceWorker)
}
