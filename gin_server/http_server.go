package gin_server

import (
	"net/http"
	"time"
)

const (
	defaultSec               = 30
	defaultReadHeaderTimeout = 3
)

type srvConf struct {
	ProxyPath         string `json:"ProxyPath"`
	ReadTimeout       int64  `json:"ReadTimeout"`
	WriteTimeout      int64  `json:"WriteTimeout"`
	IdleTimeoutx      int64  `json:"IdleTimeoutx"`
	ReadHeaderTimeout int64  `json:"ReadHeaderTimeout"`
}

func createServer(router http.Handler, addr string, cfg *srvConf) *http.Server {
	srv := http.Server{
		Addr:              addr,
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeoutx) * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	return &srv
}

func validateConfigSrv(cfg *srvConf) {
	if cfg.IdleTimeoutx <= 0 {
		cfg.IdleTimeoutx = defaultSec
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = defaultSec
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = defaultSec
	}
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = defaultReadHeaderTimeout
	}
}
