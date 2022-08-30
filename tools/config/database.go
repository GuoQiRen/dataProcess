package config

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

func InitDatabase(cfg *viper.Viper) *Database {
	return &Database{
		Dbtype:   cfg.GetString("dbtype"),
		Environ:  cfg.GetString("environ"),
		Port:     cfg.GetInt("port"),
		Host:     cfg.GetString("host"),
		Name:     cfg.GetString("name"),
		Username: cfg.GetString("username"),
		Password: cfg.GetString("password"),
	}
}
