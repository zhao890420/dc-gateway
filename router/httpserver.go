package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhao890420/dc-gateway/common"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(common.GetConfig().MustValue("base", "debug_mode", ""))
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           common.GetConfig().MustValue("http", "addr", ""),
		Handler:        r,
		ReadTimeout:    time.Duration(common.GetConfig().MustInt("http", "read_timeout", 10)) * time.Second,
		WriteTimeout:   time.Duration(common.GetConfig().MustInt("http", "write_timeout", 10)) * time.Second,
		MaxHeaderBytes: 1 << uint(common.GetConfig().MustInt("http", "max_header_bytes", 10)),
	}
	go func() {
		log.Printf(" [INFO] HttpServerRun:%s\n", common.GetConfig().MustValue("http", "addr", ""))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", common.GetConfig().MustValue("http", "addr", ""), err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
