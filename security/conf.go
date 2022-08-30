package security

import (
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

var TkInfo *TokInfo

type TokInfo struct {
	User        string `json:"userAccount"`
	Ver         int32  `json:"ver"`
	Salt        string `json:"salt"`
	Expire      int64  `json:"expire"`
	UrlUserInfo string `json:"urlUserInfo"`
	Token       string `json:"token"`
}

func GetTokInfo(path string) (err error) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	//Replace environment variables
	reader := strings.NewReader(string(content))
	err = viper.ReadConfig(reader)
	if err != nil {
		return
	}

	// 初始化
	if TkInfo == nil {
		TkInfo = new(TokInfo)
	}

	TkInfo.User = viper.GetString("token.user")
	TkInfo.Ver = viper.GetInt32("token.ver")
	TkInfo.Salt = viper.GetString("token.salt")
	TkInfo.Expire = viper.GetInt64("token.expire")
	TkInfo.UrlUserInfo = viper.GetString("token.urlUserInfo")

	return
}

func (t *TokInfo) Write() (err error) {
	viper.Set("token.ver", t.Ver)
	viper.Set("token.salt", string(t.Salt))
	writeConf := viper.WriteConfig()

	if writeConf != nil && len(writeConf.Error()) != 0 {
		err = errors.New(writeConf.Error())
		return
	}

	return
}

type UrlInfo struct {
	UrlUserInfo string `json:"urlUserInfo"`
}

func Misc() UrlInfo {
	return UrlInfo{UrlUserInfo: TkInfo.UrlUserInfo}
}
