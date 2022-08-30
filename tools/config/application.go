package config

import "github.com/spf13/viper"

type Application struct {
	ReadTimeout   int
	WriterTimeout int
	Host          string
	Port          string
	Name          string
	JwtSecret     string
	Mode          string
	DemoMsg       string
	Domain        string
	IsHttps       bool
	SrcSource     string
}

func InitApplication(cfg *viper.Viper) *Application {
	return &Application{
		ReadTimeout:   cfg.GetInt("readTimeout"),
		WriterTimeout: cfg.GetInt("writerTimeout"),
		Host:          cfg.GetString("host"),
		Port:          cfg.GetString("port"),
		Name:          cfg.GetString("name"),
		JwtSecret:     cfg.GetString("jwtSecret"),
		Mode:          cfg.GetString("mode"),
		DemoMsg:       cfg.GetString("demoMsg"),
		Domain:        cfg.GetString("domain"),
		IsHttps:       cfg.GetBool("ishttps"),
		SrcSource:     cfg.GetString("srcSource"),
	}
}
