package cmd

import (
	"context"
	config2 "dataProcess/config"
	"dataProcess/database"
	"dataProcess/logger"
	"dataProcess/router"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	config   string
	port     string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "processing server config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration filej")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8008", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func usage() {
	usageStr := `starting api server`
	log.Printf("%s\n", usageStr)
}

func setup() {
	// 读取配置文件
	config2.ConfigureSetUp(config)
	// 初始化数据库
	database.Setup()
}

func run() error {

	// 初始化路由
	r := router.InitRouter()

	srv := &http.Server{
		Addr:    config2.ApplicationConfig.Host + ":" + config2.ApplicationConfig.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if config2.ApplicationConfig.IsHttps {
			if err := srv.ListenAndServeTLS(config2.SslConfig.Pem, config2.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("listen: %s\n", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("listen: %s\n", err)
			}
		}
	}()
	fmt.Printf("Server Run http://%s:%s/ \r\n",
		config2.ApplicationConfig.Host,
		config2.ApplicationConfig.Port)
	fmt.Printf("Swagger URL http://%s:%s/swagger/index.html \r\n",
		config2.ApplicationConfig.Host,
		config2.ApplicationConfig.Port)
	fmt.Printf("Enter Control + C Shutdown Server \r\n")
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("Shutdown Server ... \r\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Info("Server exiting")
	return nil
}
