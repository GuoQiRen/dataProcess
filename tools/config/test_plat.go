package config

import "github.com/spf13/viper"

type TestConf struct {
	Head              string
	Host              string
	Port              string
	LoginContext      string
	DataSetContext    string
	UserGroupContext  string
	UserTokenContext  string
	UserDetailContext string
	UserDetailInfo    string
}

func InitTestConf(cfg *viper.Viper) *TestConf {
	return &TestConf{
		Head:              cfg.GetString("head"),
		Host:              cfg.GetString("host"),
		Port:              cfg.GetString("port"),
		LoginContext:      cfg.GetString("loginContext"),
		DataSetContext:    cfg.GetString("dataSetContext"),
		UserGroupContext:  cfg.GetString("userGroupContext"),
		UserTokenContext:  cfg.GetString("userTokenContext"),
		UserDetailContext: cfg.GetString("userDetailContext"),
		UserDetailInfo:    cfg.GetString("userDetailInfo"),
	}
}
