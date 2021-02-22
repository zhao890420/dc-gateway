package server

import (
	"github.com/spf13/cobra"
	"github.com/zhao890420/dc-gateway/common"
	"github.com/zhao890420/dc-gateway/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	config   string
	port     string
	apoEnv   string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "server",
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

/**
 * @Author zhaoguang
 * @Description 启动管理后台服务
 * @Date 10:57 下午 2021/2/22
 **/
func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "conf/config.ini", "Start server with provided configuration file")
}

func usage() {
	usageStr := `starting api server`
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
	// 5.初始化redis
	common.InitRedis()

}

func run() error {
	router.HttpServerRun()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	router.HttpServerStop()
	for _, component := range common.Destroyables {
		component.Destroy()
	}
	return nil
}
