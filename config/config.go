package config

import (
	"dataProcess/logger"
	"dataProcess/tools/config"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var (
	cfgDatabase, cfgCache, cfgApplication, cfgSsl, cfgJinn, cfgTrainPlat, cfgNginx, cfgTestPlat *viper.Viper
)

var (
	DatabaseConfig    = new(config.Database)
	CacheConfig       = new(config.Cache)
	ApplicationConfig = new(config.Application)
	SslConfig         = new(config.Ssl)
	JinnConfig        = new(config.JinnDownload)
	TrainPlatConfig   = new(config.TrainConf)
	TestPlatConfig    = new(config.TestConf)
	NginxConfig       = new(config.NginxInfo)
)

// 载入配置文件
func ConfigureSetUp(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Read config filej fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Parse config filej fail: %s", err.Error()))
	}

	// 数据库初始化
	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("config not found settings.database")
	}
	DatabaseConfig = config.InitDatabase(cfgDatabase)

	// 初始化redis缓存参数
	cfgCache = viper.Sub("settings.redis")
	if cfgCache == nil {
		panic("config not found settings.cache")
	}
	CacheConfig = config.InitCache(cfgCache)

	// 启动参数
	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = config.InitApplication(cfgApplication)

	// ssl 配置
	cfgSsl = viper.Sub("settings.ssl")
	if cfgSsl == nil {
		panic("config not found settings.ssl")
	}
	SslConfig = config.InitSsl(cfgSsl)

	// remote 配置
	cfgJinn = viper.Sub("settings.jinn")
	if cfgJinn == nil {
		panic("config not found settings.jinn")
	}
	JinnConfig = config.InitJinnDownload(cfgJinn)

	// train_plat 配置
	cfgTrainPlat = viper.Sub("settings.trainPlat")
	if cfgTrainPlat == nil {
		panic("config not found settings.trainPlat")
	}
	TrainPlatConfig = config.InitTrainConf(cfgTrainPlat)

	// test_plat 配置
	cfgTestPlat = viper.Sub("settings.testPlat")
	if cfgTestPlat == nil {
		panic("config not found settings.testPlat")
	}
	TestPlatConfig = config.InitTestConf(cfgTestPlat)

	// method_desc配置
	cfgNginx = viper.Sub("settings.nginx")
	if cfgNginx == nil {
		panic("config not found settings.nginx")
	}
	NginxConfig = config.InitNginxInfo(cfgNginx)

	// 日志配置
	logger.Init()
}
