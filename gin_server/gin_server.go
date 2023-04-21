package gin_server

import (
	"context"
	"github.com/achekslime/core/app"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartGin(engine *gin.Engine, addr string) app.ServerStartTask {
	return func(ctx context.Context) {
		startGin(ctx, engine, addr)
	}
}

func startGin(ctx context.Context, engine *gin.Engine, addr string) {
	// получить gin_config
	// для этого нужно передавать config_manager

	//cfg := srvConf{}
	//if conf != nil {
	//	ok, res := conf.GetByServiceKey("gin_config")
	//	if ok {
	//		err := jsoniter.Unmarshal(res, &cfg)
	//		if err != nil {
	//			log.Println("cant unmarshal gin conf.\nError:", err)
	//		}
	//	}
	//}

	cfg := srvConf{}
	validateConfigSrv(&cfg)

	httpServer := createServer(engine, addr, &cfg)

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Error("Server forced to shutdown:", err)
		}
	}()
}
