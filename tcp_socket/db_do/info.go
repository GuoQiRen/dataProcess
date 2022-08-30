package dbDo

import "github.com/spf13/viper"

type Database struct {
	Dbtype   string
	Environ  string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

func InitDatabaseConf(cfg *viper.Viper) *Database {
	return &Database{
		Dbtype:   cfg.GetString("settings.database.dbtype"),
		Environ:  cfg.GetString("settings.database.environ"),
		Port:     cfg.GetInt("settings.database.port"),
		Host:     cfg.GetString("settings.database.host"),
		Name:     cfg.GetString("settings.database.name"),
		Username: cfg.GetString("settings.database.username"),
		Password: cfg.GetString("settings.database.password"),
	}
}
