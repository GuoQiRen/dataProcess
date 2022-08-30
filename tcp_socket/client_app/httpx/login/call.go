package login

import (
	"dataProcess/cache/redis"
	"dataProcess/constants"
	"dataProcess/constants/train"
	"dataProcess/repository/mocks"
	"dataProcess/tcp_socket/client_app/httpx"
	dbDo "dataProcess/tcp_socket/db_do"
	"dataProcess/tools/utils/utils"
	"errors"
)

func GetLoginHead(createId int32) (head map[string]string, err error) {
	headMap := make(map[string]string)
	headMap[train.ContentKey] = train.ContentValue

	token, err := redis.GetOperate(constants.JinnToken + constants.SubLine + utils.IntToString(int(createId)))
	if err != nil {
		return
	}

	// 加入登录的token
	headMap[train.CookieKey] = train.JinnCookiePreValue + token
	return headMap, err
}

func LoginCall(headMap map[string]string) (loginResp mocks.LoginResp, err error) {
	encryptBd := httpx.CreateEncryptionBody(httpx.SetEncodeText(dbDo.DataSetConfig.Password))
	encryptPasswd, err := encryptBd.Encrypt()
	if err != nil {
		return
	}

	loginBody := httpx.CreateLoginReq(dbDo.DataSetConfig.Username, encryptPasswd, dbDo.DataSetConfig.Domain)
	resp, err := loginBody.RequestRemote(nil, headMap)
	if err != nil {
		return
	}

	// assert -> mocks.LoginResp
	loginResp = resp.(mocks.LoginResp)

	if !loginResp.Result {
		err = errors.New(loginResp.ErrorMsg)
		return
	}
	return
}

func RefreshToken(token string, headMap map[string]string) (newToken string, err error) {
	resp, err := httpx.CreateOldTokenContext(token).RequestRemote(nil, headMap)
	if err != nil {
		return
	}

	// assert -> mocks.TokenUpdateContext
	tokenBodyResp := resp.(mocks.TokenUpdateContext)

	if !tokenBodyResp.Result {
		err = errors.New(tokenBodyResp.ErrorMsg)
		return
	}

	newToken = tokenBodyResp.Token
	return
}
