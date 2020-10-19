package gateway

import (
	"encoding/json"
	"github.com/seanbit/gokit/validate"
	"io/ioutil"
	"reflect"
	"regexp"
	"unicode"
)

var _data_types_ = make(map[string]reflect.Type)

const (
	type_bool 		= "bool"
	type_int 		= "int"
	type_int8 		= "int8"
	type_int16		= "int16"
	type_int32		= "int32"
	type_int64		= "int64"
	type_uint		= "uint"
	type_uint8		= "uint8"
	type_uint16		= "uint16"
	type_uint32		= "uint32"
	type_uint64		= "uint64"
	type_float32	= "float32"
	type_float64	= "float64"
	type_string		= "string"
)

type dtInfo struct {
	Name string `json:"_data_type_name_" validate:"required,gte=1"`
	FieldType map[string]string `json:"_dt_field_type_" validate:"gt=0,dive,required"`
	FieldTag map[string]string `json:"_dt_field_tag_" validate:"gt=0,dive,required"`
}

var (
	ftRegexp  = regexp.MustCompile(`^\{([\s\S]+)}$`)
	ftsRegexp = regexp.MustCompile(`^\[([\s\S]+)]$`)
	dtsRegexp = regexp.MustCompile(`^\[]([\s\S]+)`)
)

func DataDefines(path string) map[string]reflect.Type {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	var datatypes = make([]dtInfo, 1)
	if err := json.Unmarshal(bts, &datatypes); err != nil {
		log.Panic(err)
	}
	for _, dtinfo := range datatypes {
		dataType := parseStructDataType(dtinfo)
		_data_types_[dtinfo.Name] = dataType
	}
	return _data_types_
}

func NewData(dtname string) interface{} {
	if dt, ok := parseBaseType(dtname); ok {
		return reflect.New(dt).Interface()
	}
	if params := dtsRegexp.FindStringSubmatch(dtname); len(params) == 2 && params[0] == dtname { // slice shoud handle
		typeName := params[1]
		if dt, ok := parseBaseType(typeName); ok {
			return reflect.MakeSlice(dt, 0, 0).Interface()
		} else if dt, ok := _data_types_[typeName]; ok {
			return reflect.MakeSlice(dt, 0, 0).Interface()
		} else {
			return nil
		}
	}
	if dt, ok := _data_types_[dtname]; ok {
		return reflect.New(dt).Interface()
	} else {
		return nil
	}
}

func parseStructDataType(info dtInfo) reflect.Type {
	if err := validate.ValidateParameter(info); err != nil {
		log.Fatalf("datainfo parse faiiled on {%s}, valid err: %s", info.Name, err.Error())
	}
	var sfs []reflect.StructField
	for field, typeInfo := range info.FieldType {
		var fieldType reflect.Type
		if ft, ok := parseBaseType(typeInfo); ok {
			fieldType = ft
		} else if params := ftRegexp.FindStringSubmatch(typeInfo); len(params) == 2 && params[0] == typeInfo {
			typeName := params[1]
			if dt, ok := _data_types_[typeName]; ok {
				fieldType = dt
			} else {
				log.Fatalf("data parse type err: could not found typename {%s} before this data {%s}", typeName, info.Name)
			}
		} else if params = ftsRegexp.FindStringSubmatch(typeInfo); len(params) == 2 && params[0] == typeInfo { // slice shoud handle
			typeName := params[1]
			if dt, ok := parseBaseType(typeName); ok {
				fieldType = reflect.SliceOf(dt)
			} else if dt, ok := _data_types_[typeName]; ok {
				fieldType = reflect.SliceOf(dt)
			} else {
				log.Fatalf("data parse type err: could not found typename [%s] before this data {%s}", typeName, info.Name)
			}
		} else {
			log.Fatalf("data parse type err: could not parse type with %s in this data %s", typeInfo, info.Name)
		}
		var fieldTag ,ok = info.FieldTag[field]
		if !ok {
			fieldTag = field
			if fcIsUpper(fieldTag) {
				fieldTag = fcLower(field)
			}
		}
		sf := reflect.StructField{
			Name: field,
			Type: fieldType,
			Tag: reflect.StructTag(fieldTag),
		}
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	return st
}

func parseBaseType(v string) (reflect.Type, bool) {
	switch v {
	case type_bool:
		return reflect.TypeOf(bool(false)), true
	case type_int:
		return reflect.TypeOf(int(0)), true
	case type_int8:
		return reflect.TypeOf(int8(0)), true
	case type_int16:
		return reflect.TypeOf(int16(0)), true
	case type_int32:
		return reflect.TypeOf(int32(0)), true
	case type_int64:
		return reflect.TypeOf(int64(0)), true
	case type_uint:
		return reflect.TypeOf(uint(0)), true
	case type_uint8:
		return reflect.TypeOf(uint8(0)), true
	case type_uint16:
		return reflect.TypeOf(uint16(0)), true
	case type_uint32:
		return reflect.TypeOf(uint32(0)), true
	case type_uint64:
		return reflect.TypeOf(uint64(0)), true
	case type_float32:
		return reflect.TypeOf(float32(0)), true
	case type_float64:
		return reflect.TypeOf(float64(0)), true
	case type_string:
		return reflect.TypeOf(string("")), true
	default:
		return nil, false
	}
}

func fcIsUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

func fcIsLower(s string) bool {
	return unicode.IsLower([]rune(s)[0])
}

func fcUpper(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func fcLower(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//func newData(dt map[string]interface{}) interface{} {
//	var sfs []reflect.StructField
//	for k, v := range dt {
//		t := reflect.TypeOf(v)
//		sf := reflect.StructField{
//			Name: k,
//			Type: t,
//		}
//		sfs = append(sfs, sf)
//	}
//	st := reflect.StructOf(sfs)
//	so := reflect.New(st)
//	return so.Interface()
//}





