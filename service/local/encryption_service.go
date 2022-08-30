package local

type EncryptServer interface {
	Encrypt() (encodeString string, err error)
	Decrypt() (decodeString string, err error)
}
