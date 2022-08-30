package linkDo

import "github.com/spf13/viper"

type LinkConf struct {
	Host    string
	Port    string
	Network string
}

func InitLinkConf(cfg *viper.Viper) *LinkConf {
	return &LinkConf{
		Host:    cfg.GetString("settings.serverConf.host"),
		Port:    cfg.GetString("settings.serverConf.port"),
		Network: cfg.GetString("settings.serverConf.network"),
	}
}
