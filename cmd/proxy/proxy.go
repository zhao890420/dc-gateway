package api

import (
	"github.com/spf13/cobra"
	"github.com/zhao890420/dc-gateway/common"
	"github.com/zhao890420/dc-gateway/dao"
	"github.com/zhao890420/dc-gateway/grpc_proxy_router"
	"github.com/zhao890420/dc-gateway/http_proxy_router"
	"github.com/zhao890420/dc-gateway/tcp_proxy_router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/**
 * @Author zhaoguang
 * @Description 启动代理服务
 * @Date 10:53 下午 2021/2/21
 **/
var (
	config   string
	port     string
	apoEnv   string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "proxy",
		Short:   "Start proxy server",
		Example: "proxy",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

const (
	ModeProd = "prod"
	ModeDev  = "dev"
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "conf/config.ini", "Start server with provided configuration file")
}

func usage() {
	usageStr := `starting proxy server`
	log.Printf(usageStr)
}

func setup() {
	// 1. 初始化日志
	common.InitLogs()
	// 2. 读取配置
	common.InitConfigFile(config)
	// 3. 初始化数据库
	common.InitDatabase()
	// 4. 初始化阿波罗配置
	//common.InitApollo(apoEnv)

}

func run() error {
	dao.ServiceManagerHandler.LoadOnce()
	dao.AppManagerHandler.LoadOnce()

	go func() {
		http_proxy_router.HttpServerRun()
	}()
	go func() {
		http_proxy_router.HttpsServerRun()
	}()
	go func() {
		tcp_proxy_router.TcpServerRun()
	}()
	go func() {
		grpc_proxy_router.GrpcServerRun()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	tcp_proxy_router.TcpServerStop()
	grpc_proxy_router.GrpcServerStop()
	http_proxy_router.HttpServerStop()
	http_proxy_router.HttpsServerStop()
	return nil
}
