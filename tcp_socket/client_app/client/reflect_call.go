package clientApp

import (
	"dataProcess/constants"
	"reflect"
)

type reflectBody struct {
	methodName    string
	methodAddress interface{}
	methodParams  []reflect.Value
	processArgs   map[string]interface{}
	mainEngine    string
	methodType    int32
}

func CreateReflectBody(methodName string, methodAddress interface{}, methodParams []reflect.Value,
	processArgs map[string]interface{}, mainEngine string, methodType int32) *reflectBody {
	return &reflectBody{methodName: methodName, methodAddress: methodAddress, methodParams: methodParams,
		processArgs: processArgs, mainEngine: mainEngine, methodType: methodType}
}

/**
反射获取方法，并调用方法
*/
func (r *reflectBody) reflectMethodCall(jsonFile, jinnToken string) []reflect.Value {
	if r.processArgs == nil {
		r.processArgs = make(map[string]interface{})
	}
	r.processArgs[constants.JsonFile] = jsonFile
	r.processArgs[constants.JinnToken] = jinnToken

	if r.methodParams == nil {
		r.methodParams = make([]reflect.Value, 3)
	}
	r.methodParams[0] = reflect.ValueOf(r.mainEngine)
	r.methodParams[1] = reflect.ValueOf(r.processArgs)
	r.methodParams[2] = reflect.ValueOf(r.methodType)

	v := reflect.ValueOf(r.methodAddress)
	return v.MethodByName(r.methodName).Call(r.methodParams)
}
