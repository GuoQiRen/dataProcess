package config

import "github.com/spf13/viper"

type ContainerConf struct {
	Host          string
	Port          string
	Network       string
	Command       string
	ContainerPort string
	UserName      string
}

func InitContainerConf(cfg *viper.Viper) *ContainerConf {
	return &ContainerConf{
		Host:          cfg.GetString("host"),
		Port:          cfg.GetString("port"),
		Network:       cfg.GetString("network"),
		ContainerPort: cfg.GetString("containerPort"),
		UserName:      cfg.GetString("userName"),
		Command:       cfg.GetString("command"),
	}
}
