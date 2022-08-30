package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"dataProcess/jinn/src/utils/jsonx"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type SessData struct {
	Token    string
	UserData TokenInfo
}

type TokenInfo struct {
	User  string `json:"userAccount"`
	Ver   int32  `json:"ver"`
	Exp   int64  `json:"exp"`
	Token string
}

func (t *TokenInfo) MarshalJson(doc []byte) []byte {
	doc = jsonx.AppendDocumentStart(doc)
	doc = jsonx.AppendStringElement(doc, "userAccount", t.User)
	doc = jsonx.AppendInt32Element(doc, "ver", t.Ver)
	doc = jsonx.AppendStringElement(doc, "platformCode", "JinnData_material")
	doc = jsonx.AppendInt64Element(doc, "exp", time.Now().Unix()+TkInfo.Expire)
	return jsonx.AppendDocumentEnd(doc).Bytes()
}

/**
获取token信息
*/
func GetTokenInfo(d *SessData) error {
	token := d.Token
	info := TokenInfo{Token: token}

	err := GetTokInfo("security/token.yaml")
	if err != nil {
		return err
	}
RETRY:
	raw, err := ParseJWT([]byte(token), []byte(TkInfo.Salt))
	if err == nil {
		if err = jsonx.Unmarshal(raw, &info); err != nil {
			return err
		}
		remainTime := info.Exp - time.Now().Unix()
		if remainTime < 0 {
			return errors.New("token expired")
		} else if remainTime < 30*60 {
			token = string(CreateJWT(info.MarshalJson(make([]byte, 0, 64)), []byte(TkInfo.Salt)))
		}
		d.UserData = info
		d.UserData.Token = token
		return nil
	} else {
		if len(raw) > 0 {
			if err = jsonx.Unmarshal(raw, &info); err != nil {
				return err
			}
			if info.Ver > TkInfo.Ver {
				ver, salt, err := UpdateSalt(token)
				if err != nil {
					return err
				}
				TkInfo.Ver = ver
				TkInfo.Salt = salt
				err = TkInfo.Write()
				if err != nil {
					return err
				}
				goto RETRY
			} else if info.Ver < TkInfo.Ver {
				token, err := UpdateToken(info)
				if err != nil {
					return err
				}
				d.UserData.Token = token
				return nil
			} else {
				return errors.New("invalid token")
			}
		}
		return err
	}
}

/**
设置cookies
*/
func SetCookie(c *gin.Context, newToken string) (err error) {
	token, err := c.Cookie("token")
	if err != nil {
		return
	}
	if len(token) == 0 {
		err = errors.New("token is nil")
		return
	}
	c.Request.Header.Set("Set_Cookie", "token="+newToken)
	return nil
}

/**
更新salt
*/
func UpdateSalt(token string) (int32, string, error) {
	req, err := http.NewRequest("GET", Misc().UrlUserInfo+"/usercenter/api/v2/update/salt", nil)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	resp, err := (&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}).Do(req)
	if err != nil {
		return 0, "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return 0, "", err
	}
	tokenSet := struct {
		Result    bool
		ErrorMsg  string
		ErrorCode int32
		Data      struct {
			Ver  int32
			Salt string
		}
	}{}
	err = jsonx.Unmarshal(body, &tokenSet)
	if err != nil {
		return 0, "", err
	}
	if !tokenSet.Result {
		return 0, "", errors.New(tokenSet.ErrorMsg)
	}
	return tokenSet.Data.Ver, tokenSet.Data.Salt, nil
}

/**
更新token
*/
func UpdateToken(info TokenInfo) (string, error) {
	doc := jsonx.AppendDocumentStart(make([]byte, 0, 64))
	doc = jsonx.AppendStringElement(doc, "token", info.Token)
	doc = jsonx.AppendDocumentEnd(doc).Bytes()
	req, err := http.NewRequest("POST", Misc().UrlUserInfo+"/usercenter/api/v2/update/token", bytes.NewReader(doc))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+string(CreateJWT(info.MarshalJson(make([]byte, 0, 64)), []byte(TkInfo.Salt))))
	resp, err := (&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}).Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return "", err
	}
	tokenSet := struct {
		Result    bool
		ErrorMsg  string
		ErrorCode int32
		Data      struct {
			Token string
		}
	}{}
	err = jsonx.Unmarshal(body, &tokenSet)
	if err != nil {
		return "", err
	}
	if !tokenSet.Result {
		return "", errors.New(tokenSet.ErrorMsg)
	}
	return tokenSet.Data.Token, nil
}

/**
创建jwt
*/
func CreateJWT(data []byte, key []byte) []byte {
	buf := make([]byte, 32+36+base64.RawURLEncoding.EncodedLen(len(data))+43+2)
	token := buf[32:]
	temp := append(token[:0], "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9."...)
	base64.RawURLEncoding.Encode(token[len(temp):], data)
	temp = token[:len(temp)+base64.RawURLEncoding.EncodedLen(len(data))]
	h := hmac.New(sha256.New, key)
	h.Write(temp)
	temp = append(temp, '.')
	base64.RawURLEncoding.Encode(token[len(temp):], h.Sum(buf[:0]))
	return token
}

/**
解析jwt
*/
func ParseJWT(token []byte, key []byte) ([]byte, error) {
	buf := make([]byte, len(token), len(token)+56)
	pl := bytes.IndexByte(token, '.') + 1
	sig := bytes.LastIndexByte(token, '.') + 1
	if pl <= 2 || sig <= pl+3 {
		return nil, errors.New("invalid token")
	}
	n, err := base64.RawURLEncoding.Decode(buf[pl:], token[pl:sig-1])
	if err != nil {
		return nil, err
	}
	if key != nil {
		h := hmac.New(sha256.New, key)
		h.Write(token[:sig-1])
		base64.RawURLEncoding.Encode(buf[sig:], h.Sum(buf[len(buf):]))
		if !bytes.Equal(buf[sig:], token[sig:]) {
			return buf[pl : pl+n], errors.New("invalid signature")
		}
	}
	return buf[pl : pl+n], nil
}
