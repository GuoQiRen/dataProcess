package clientApp

import (
	"dataProcess/constants"
	"dataProcess/constants/method"
	"dataProcess/constants/plat"
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type MSummary struct {
}

func (m *MSummary) MethodExec(mainEngine string, args map[string]interface{}, methodType int32) (pro *os.Process, errMsg string) {
	params := m.transformToString(args)

	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, //其他变量如果不清楚可以不设定
	}

	mainArgs := m.transformToLinkParams(mainEngine, params, methodType)
	engineDir := getEngineDir(mainEngine)
	err := os.Chdir(engineDir)
	if err != nil {
		return
	}
	pro, err = os.StartProcess(casePlatCmd(methodType), mainArgs, attr) // 开始进程
	if err != nil {
		errMsg = err.Error()
		return
	}

	return
}

func getEngineDir(mainEngine string) (engineDir string) {
	mainEngineArr := strings.Split(mainEngine, string(os.PathSeparator))
	engineArr := mainEngineArr[:len(mainEngineArr)-1]
	engineDir = strings.Join(engineArr, string(os.PathSeparator))
	return
}

func casePlatCmd(methodType int32) (cmd string) {
	if methodType == method.PubMethod {
		switch runtime.GOOS {
		case plat.Windows:
			cmd = plat.WindowsPython3
		case plat.Linux:
			cmd = plat.LinuxPython3
		}
		return cmd
	}

	if methodType == method.MyMethod || methodType == method.ShareMethod {
		cmd = plat.Sh
	}
	return
}

func (m *MSummary) transformToString(args map[string]interface{}) (params map[string]string) {
	var tran string

	params = make(map[string]string)
	for key, val := range args {
		item := reflect.ValueOf(val)
		paramType := item.Kind()

		switch paramType {
		case reflect.String:
			tran = item.Interface().(string)
		case reflect.Int:
			tran = strconv.Itoa(item.Interface().(int))
		case reflect.Int32:
			tran = strconv.FormatInt(int64(item.Interface().(int32)), 10)
		case reflect.Int64:
			tran = strconv.FormatInt(item.Interface().(int64), 10)
		case reflect.Bool:
			tran = strconv.FormatBool(item.Interface().(bool))
		case reflect.Float32:
			tran = strconv.FormatFloat(float64(item.Interface().(float32)), 'f', 10, 64)
		case reflect.Float64:
			tran = strconv.FormatFloat(item.Interface().(float64), 'f', 10, 64)
		default:
			inter, _ := json.Marshal(item.Interface())
			tran = string(inter)
		}
		params[key] = tran
	}
	return
}

func (m *MSummary) transformToLinkParams(mainEngine string, params map[string]string, methodType int32) (args []string) {
	args = append(args, casePlatCmd(methodType), mainEngine)

	for key, val := range params {
		if key == "" {
			args = append(args, val)
		}

		args = append(args, constants.ParamsPrefix+key, val)
	}

	return
}
