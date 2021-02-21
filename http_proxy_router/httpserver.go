package http_proxy_router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhao890420/dc-gateway/common"
	"github.com/zhao890420/dc-gateway/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(common.GetConfig().MustValue("base", "debug_mode", ""))
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           common.GetConfig().MustValue("proxy_http", "addr", ""),
		Handler:        r,
		ReadTimeout:    time.Duration(common.GetConfig().MustInt("proxy_http", "read_timeout", 10)) * time.Second,
		WriteTimeout:   time.Duration(common.GetConfig().MustInt("proxy_http", "write_timeout", 10)) * time.Second,
		MaxHeaderBytes: 1 << uint(common.GetConfig().MustInt("proxy_http", "max_header_bytes", 10)),
	}
	log.Printf(" [INFO] http_proxy_run %s\n", common.GetConfig().MustValue("proxy_http", "addr", ""))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] http_proxy_run %s err:%v\n", common.GetConfig().MustValue("proxy_http", "addr", ""), err)
	}
}

func HttpsServerRun() {
	gin.SetMode(common.GetConfig().MustValue("base", "debug_mode", ""))
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           common.GetConfig().MustValue("proxy_https", "addr", ""),
		Handler:        r,
		ReadTimeout:    time.Duration(common.GetConfig().MustInt("proxy_https", "read_timeout", 10)) * time.Second,
		WriteTimeout:   time.Duration(common.GetConfig().MustInt("proxy_https", "write_timeout", 10)) * time.Second,
		MaxHeaderBytes: 1 << uint(common.GetConfig().MustInt("proxy_https", "max_header_bytes", 10)),
	}
	log.Printf(" [INFO] https_proxy_run %s\n", common.GetConfig().MustValue("proxy_https", "addr", ""))
	//todo 以下命令只在编译机有效，如果是交叉编译情况下需要单独设置路径
	//if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err!=http.ErrServerClosed {
	if err := HttpsSrvHandler.ListenAndServeTLS("./cert_file/server.crt", "./cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] https_proxy_run %s err:%v\n", common.GetConfig().MustValue("proxy_https", "addr", ""), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop %v stopped\n", common.GetConfig().MustValue("proxy_http", "addr", ""))
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop %v stopped\n", common.GetConfig().MustValue("proxy_https", "addr", ""))
}
