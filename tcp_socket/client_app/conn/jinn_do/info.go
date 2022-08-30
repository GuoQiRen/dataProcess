package jinnDo

import "github.com/spf13/viper"

type JinnDo struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Domain    string
	BatchSize int
}

func InitJinnDo(cfg *viper.Viper) *JinnDo {
	return &JinnDo{
		Host:      cfg.GetString("settings.jinn.host"),
		Port:      cfg.GetInt("settings.jinn.port"),
		Username:  cfg.GetString("settings.jinn.username"),
		Password:  cfg.GetString("settings.jinn.password"),
		Domain:    cfg.GetString("settings.jinn.domain"),
		BatchSize: cfg.GetInt("settings.jinn.batchsize"),
	}
}
