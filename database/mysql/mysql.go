package mysql

import (
	"bytes"
	"dataProcess/config"
	"dataProcess/logger"
	"dataProcess/orm"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"strconv"
)

var DbType string

type Mysql struct {
}

func (m *Mysql) SetUp() {
	var err error

	db := new(Mysql)
	orm.MysqlConn = db.GetConnect()
	orm.MysqlDb, err = db.Open(DbType, orm.MysqlConn)

	if err != nil {
		logger.Fatalf("%s connect error %v", DbType, err)
	} else {
		logger.Infof("%s connect success!", DbType)
	}

	if orm.MysqlDb.Error != nil {
		logger.Fatalf("serverDb error %v", orm.MysqlDb.Error)
	}

	// 是否开启详细日志记录
	orm.MysqlDb.LogMode(viper.GetBool("settings.gorm.logMode"))

	// 设置最大打开连接数
	orm.MysqlDb.DB().SetMaxOpenConns(viper.GetInt("settings.gorm.maxOpenConn"))

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	orm.MysqlDb.DB().SetMaxIdleConns(viper.GetInt("settings.gorm.maxIdleConn"))
}

func (m *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

func (m *Mysql) GetConnect() string {

	DbType = config.DatabaseConfig.Dbtype
	Host := config.DatabaseConfig.Host
	Port := config.DatabaseConfig.Port
	Name := config.DatabaseConfig.Name
	Username := config.DatabaseConfig.Username
	Password := config.DatabaseConfig.Password

	var conn bytes.Buffer
	conn.WriteString(Username)
	conn.WriteString(":")
	conn.WriteString(Password)
	conn.WriteString("@tcp(")
	conn.WriteString(Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(Port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(Name)
	conn.WriteString("?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms")
	return conn.String()
}
