package datasetDo

import "github.com/spf13/viper"

type DataSetDo struct {
	Head                       string
	Host                       string
	Port                       string
	Username                   string
	Password                   string
	Domain                     string
	LoginContext               string
	UserTokenContext           string
	PreCreateDatasetContext    string
	FormalCreateDatasetContext string
}

func InitDataSetDo(cfg *viper.Viper) *DataSetDo {
	return &DataSetDo{
		Head:                       cfg.GetString("settings.dataSet.head"),
		Host:                       cfg.GetString("settings.dataSet.host"),
		Port:                       cfg.GetString("settings.dataSet.port"),
		Username:                   cfg.GetString("settings.dataSet.username"),
		Password:                   cfg.GetString("settings.dataSet.password"),
		Domain:                     cfg.GetString("settings.dataSet.domain"),
		LoginContext:               cfg.GetString("settings.dataSet.loginContext"),
		UserTokenContext:           cfg.GetString("settings.dataSet.userTokenContext"),
		PreCreateDatasetContext:    cfg.GetString("settings.dataSet.PreCreateDatasetContext"),
		FormalCreateDatasetContext: cfg.GetString("settings.dataSet.FormalCreateDatasetContext"),
	}
}
