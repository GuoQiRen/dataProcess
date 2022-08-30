package dbDo

import (
	"bytes"
	"dataProcess/constants"
	"dataProcess/logger"
	"dataProcess/tools/utils/utils"
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
)

var DbType string

type MongodbServer struct {
}

func (m *MongodbServer) SetUp() (mongoConn *mgo.Database) {
	var err error
	db := new(MongodbServer)

	envi := DatabaseConfig.Environ

	if envi == constants.Formalization {
		mongoConn, err = db.OpenFormalization()
	} else {
		mongoConn, err = db.Open(GetConnect())
	}
	if err != nil {
		logger.Fatalf("%s connect error %v", DbType, err)
		return
	}
	return mongoConn
}

func (m *MongodbServer) Open(conn string) (db *mgo.Database, err error) {
	dbName := DatabaseConfig.Name

	session, err := mgo.Dial(conn)
	if err != nil {
		logger.Fatalf("serverDb error %v", err.Error())
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(dbName)
	return
}

func (m *MongodbServer) OpenFormalization() (db *mgo.Database, err error) {
	DbType = DatabaseConfig.Dbtype
	Host := DatabaseConfig.Host
	Port := DatabaseConfig.Port
	Username := DatabaseConfig.Username
	Password := DatabaseConfig.Password
	dbName := DatabaseConfig.Name

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
	db = session.DB(dbName)
	return
}

func GetConnect() string {

	DbType = DatabaseConfig.Dbtype
	Host := DatabaseConfig.Host
	Port := DatabaseConfig.Port
	Username := DatabaseConfig.Username
	Password := DatabaseConfig.Password

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
