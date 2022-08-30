package impl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"dataProcess/constants/encrypt"
	"encoding/base64"
	"encoding/hex"
)

type EncryptionBody struct {
	key        []byte
	iv         []byte
	encodeText string
	decodeText string
}

type Option func(*EncryptionBody)

func SetKey(key []byte) Option {
	return func(this *EncryptionBody) {
		this.key = key
	}
}

func SetIv(iv []byte) Option {
	return func(this *EncryptionBody) {
		this.iv = iv
	}
}

func SetEncodeText(encodeText string) Option {
	return func(this *EncryptionBody) {
		this.encodeText = encodeText
	}
}

func SetDecodeText(decodeText string) Option {
	return func(this *EncryptionBody) {
		this.decodeText = decodeText
	}
}

func CreateEncryptionBody(opts ...Option) EncryptionBody {
	sixthKey := md5.Sum([]byte(encrypt.DefaultKey))
	sixthIv := md5.Sum([]byte(encrypt.DefaultIv))
	twoKey := sixthKey[0:]
	twoIv := sixthIv[0:]
	key := hex.EncodeToString(twoKey)
	iv := hex.EncodeToString(twoIv)
	ivByte := []byte(iv)
	defaultEncryption := EncryptionBody{
		key: []byte(key),
		iv:  ivByte[8 : len(iv)-8],
	}

	for _, o := range opts {
		o(&defaultEncryption)
	}

	return defaultEncryption
}

func (e *EncryptionBody) Encrypt() (encodeString string, err error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return
	}

	add := 0
	length := 16
	count := len(e.encodeText)
	if count%length != 0 {
		add = length - (count % length)
	}

	zeroText := ""
	for i := 0; i < add; i++ {
		zeroText += "\x00"
	}
	plainText := []byte(e.encodeText + zeroText) // 填补的空白

	encodeCiphertext := make([]byte, len(e.encodeText+zeroText))

	mode := cipher.NewCBCEncrypter(block, e.iv)
	mode.CryptBlocks(encodeCiphertext, plainText)

	return base64.StdEncoding.EncodeToString(encodeCiphertext), err
}

func (e *EncryptionBody) Decrypt() (decodeString string, err error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return
	}
	mode := cipher.NewCBCDecrypter(block, e.iv)

	ciphertext, err := base64.StdEncoding.DecodeString(e.decodeText)
	if err != nil {
		return
	}
	decodeCiphertext := make([]byte, len(ciphertext))
	mode.CryptBlocks(decodeCiphertext, ciphertext)

	bytes.TrimRight(decodeCiphertext, "\x00")

	return string(decodeCiphertext), err
}
