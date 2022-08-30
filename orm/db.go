package orm

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

var MysqlDb *gorm.DB
var MongoDb *mgo.Database
var MysqlConn string
var MongoConn string
