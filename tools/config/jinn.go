package config

import "github.com/spf13/viper"

type JinnDownload struct {
	Host      string
	Port      int
	Username  string
	GroupName string
	Password  string
	Domain    string
	BatchSize int
}

func InitJinnDownload(cfg *viper.Viper) *JinnDownload {
	return &JinnDownload{
		Host:      cfg.GetString("host"),
		Port:      cfg.GetInt("port"),
		Username:  cfg.GetString("username"),
		GroupName: cfg.GetString("groupname"),
		Password:  cfg.GetString("password"),
		Domain:    cfg.GetString("domain"),
		BatchSize: cfg.GetInt("batchsize"),
	}
}
