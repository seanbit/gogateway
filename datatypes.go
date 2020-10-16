package gateway

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
)

var __DataTypes map[string]map[string]interface{}

func LoadDatas(path string) map[string]map[string]interface{} {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	var datatypes = make(map[string]map[string]interface{})
	if err := json.Unmarshal(bts, &datatypes); err != nil {
		log.Panic(err)
	}
	__DataTypes = datatypes
	return datatypes
}

func NewData(dtname string) interface{} {
	if data, ok := __DataTypes[dtname]; ok {
		return newData(data)
	} else {
		return nil
	}
}

func newData(dt map[string]interface{}) interface{} {
	var sfs []reflect.StructField
	for k, v := range dt {
		t := reflect.TypeOf(v)
		sf := reflect.StructField{
			Name: k,
			Type: t,
		}
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	so := reflect.New(st)
	return so.Interface()
}