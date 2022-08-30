package dbDo

import (
	"dataProcess/tcp_socket/client_app/conn/datasetDo"
	jinnDo "dataProcess/tcp_socket/client_app/conn/jinn_do"
	linkDo "dataProcess/tcp_socket/client_app/conn/link_do"
	"dataProcess/tcp_socket/client_app/conn/redis_do"
	"github.com/spf13/viper"
)

var (
	LinkConfig     = new(linkDo.LinkConf)
	DatabaseConfig = new(Database)
	JinnConfig     = new(jinnDo.JinnDo)
	DataSetConfig  = new(datasetDo.DataSetDo)
	RedisConfig    = new(redisDo.RedisDo)
)

func ConfigureSetUp(filePath string) {
	conf := viper.New()
	conf.SetConfigType("yaml")
	conf.SetConfigFile(filePath)

	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}

	// 读取server配置
	LinkConfig = linkDo.InitLinkConf(conf)

	// jinn操作配置
	JinnConfig = jinnDo.InitJinnDo(conf)

	// redis操作配置
	RedisConfig = redisDo.InitRedisDo(conf)

	// dataSet操作配置
	DataSetConfig = datasetDo.InitDataSetDo(conf)

	// database配置
	DatabaseConfig = InitDatabaseConf(conf)
}
