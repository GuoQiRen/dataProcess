package utils

import (
	"dataProcess/constants/plat"
	"errors"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SplitJinnString(path string) (space, context string, err error) {
	ret := strings.Split(path, plat.LinuxSpiltRex)
	if len(ret) < 2 {
		err = errors.New("jinnPath is invalid")
		return
	}
	space = ret[1]
	for index, s := range ret {
		if index > 1 {
			context += plat.LinuxSpiltRex + s
		}
	}
	return
}
