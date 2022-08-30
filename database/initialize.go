package database

import (
	"dataProcess/config"
	"dataProcess/database/mongo"
	"dataProcess/database/mysql"
)

func Setup() {
	dbType := config.DatabaseConfig.Dbtype
	if dbType == "mysql" {
		var db = new(mysql.Mysql)
		db.SetUp()
	} else if dbType == "mongodb" {
		var db = new(mongo.Mongodb)
		db.SetUp()
	}
}
