package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	p "github.com/zhao890420/dc-gateway/cmd/proxy"
	s "github.com/zhao890420/dc-gateway/cmd/server"
	"log"
	"os"
)

/**
 * 当前项目的可执行文件。cmd 目录下的每一个子目录名称都应该匹配可执行文件。
 * 比如果我们的项目是一个 GRpc 服务，在 /cmd/myApp/cmd.go 中就包含了启动服务进程的代码，编译后生成的可执行文件就是 myApp
 * @author zhaoguang
 * @Date 2020/10/4 10:07 上午
 */
var rootCmd = &cobra.Command{
	Use:               "gateway",
	Short:             "-v",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `gateway`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("add arg `server` to start api server, `help` to list help info !")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		usagr := `-h 查看命令`
		log.Println(usagr)
	},
}

func init() {
	rootCmd.AddCommand(s.StartCmd)
	rootCmd.AddCommand(p.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
