package mongo

import (
	"bytes"
	"dataProcess/config"
	"dataProcess/constants"
	"dataProcess/logger"
	"dataProcess/orm"
	"dataProcess/tools/utils/utils"
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
)

var DbType string

type Mongodb struct {
}

func (m *Mongodb) SetUp() {
	var err error
	db := new(Mongodb)

	envi := config.DatabaseConfig.Environ

	if envi == constants.Formalization {
		orm.MongoDb, err = db.OpenFormalization()
	} else {
		orm.MongoDb, err = db.Open(GetConnect())
	}
	if err != nil {
		logger.Fatalf("%s connect error %v", DbType, err)
	} else {
		logger.Infof("%s connect success!", DbType)
	}
}

func (m *Mongodb) Open(conn string) (db *mgo.Database, err error) {
	dbName := config.DatabaseConfig.Name

	session, err := mgo.Dial(conn)
	if err != nil {
		logger.Fatalf("serverDb error %v", err.Error())
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(dbName)
	return
}

func (m *Mongodb) OpenFormalization() (db *mgo.Database, err error) {
	Name := config.DatabaseConfig.Name
	DbType = config.DatabaseConfig.Dbtype
	Host := config.DatabaseConfig.Host
	Port := config.DatabaseConfig.Port
	Username := config.DatabaseConfig.Username
	Password := config.DatabaseConfig.Password
	dbName := config.DatabaseConfig.Name

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{Host + constants.Colon + utils.IntToString(Port)},
		Timeout:  60 * time.Second,
		Username: Username,
		Password: Password,
		Database: dbName,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Fatalf("serverDb error %v", err.Error())
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(Name)
	return
}

func GetConnect() string {

	DbType = config.DatabaseConfig.Dbtype
	Host := config.DatabaseConfig.Host
	Port := config.DatabaseConfig.Port
	Username := config.DatabaseConfig.Username
	Password := config.DatabaseConfig.Password

	var conn bytes.Buffer
	conn.WriteString(DbType)
	conn.WriteString(":")
	conn.WriteString("//")
	conn.WriteString(Username)
	conn.WriteString(":")
	conn.WriteString(Password)
	conn.WriteString("@")
	conn.WriteString(Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(Port))
	return conn.String()
}
