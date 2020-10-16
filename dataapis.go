package gateway

import (
	"encoding/json"
	"github.com/seanbit/gokit/validate"
	"io/ioutil"
)

var __APIServices []APIService

type HTTPMethod string
const (
	HTTPMethodPOST 		= "post"
	HTTPMethodGET 		= "get"
	HTTPMethodPUT 		= "put"
	HTTPMethodDELETE 	= "delete"
)

type APIService struct {
	Path        	string     		`json:"path" validate:"required,gte=1"`
	MiddelWares 	[]string   		`json:"middelWares"`
	Services 		[]APIService	`json:"services" validate:"required,gte=0,dive,required"`
	Do 				APIDo			`json:"do" validate:"required"`
}

type APIDo struct {
	HttpMethod  	HTTPMethod 		`json:"httpMethod" validate:"required,oneof=post get put delete"`
	RpcServer   	string 	   		`json:"rpcServer" validate:"required,gte=1"`
	RpcService  	string 	   		`json:"rpcService" validate:"required,gte=1"`
	RpcMethod   	string     		`json:"rpcMethod" validate:"required,gte=1"`
	RpcRequest  	string     		`json:"rpcRequest" validate:"required,gte=1"`
	RpcResponse 	string     		`json:"rpcResponse" validate:"required,gte=1"`
}

func LoadApis(path string) []APIService {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	var apiservices = make([]APIService, 1)
	if err := json.Unmarshal(bts, &apiservices); err != nil {
		log.Panic(err)
	}
	if err := validate.ValidateParameter(apiservices); err != nil {
		log.Panic(err)
	}
	__APIServices = apiservices
	return apiservices
}


