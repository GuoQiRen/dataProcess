package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func GenerateJsonFile(body interface{}, filePath string) (err error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filePath, buf, os.ModePerm)
	if err != nil {
		return
	}
	return
}

func ObtainJsonFileInfo(body interface{}, filePath string) (err error) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, body)
	if err != nil {
		return
	}
	return
}
