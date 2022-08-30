package utils

import (
	"errors"
	"strconv"
	"time"
)

func IntToString(key int) (ret string) {
	keyS := strconv.Itoa(key)
	return keyS
}

func StringToint32(key string) (ret int32, err error) {
	if len(key) == 0 {
		err = errors.New("key is empty")
		return 0, err
	}

	keyI, err := strconv.Atoi(key)
	return int32(keyI), err
}

func TimeToString(timeStamp time.Time) (timeStr string) {
	return timeStamp.Format("2006-01-02 15:04:05")
}
