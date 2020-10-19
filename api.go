package gateway

import (
	"encoding/json"
	"github.com/seanbit/gokit/validate"
	"io/ioutil"
)

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
	HttpMethod   HTTPMethod `json:"httpMethod" validate:"required,oneof=post get put delete"`
	RpcServer    string     `json:"rpcServer" validate:"required,gte=1"`
	RpcService   string     `json:"rpcService" validate:"required,gte=1"`
	RpcMethod    string     `json:"rpcMethod" validate:"required,gte=1"`
	RpcParameter string     `json:"rpcParameter" validate:"required,gte=1"`
	RpcResponse  string     `json:"rpcResponse" validate:"required,gte=1"`
}

var _api_services_ []APIService

func ApiDefines(path string) []APIService {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	type ApiServices struct {
		apiservices *[]APIService	`validate:"required,gte=0,dive,required"`
	}
	var apiservices = &ApiServices{apiservices: new([]APIService)}
	if err := json.Unmarshal(bts, apiservices.apiservices); err != nil {
		log.Panic(err)
	}
	if err := validate.ValidateParameter(apiservices); err != nil {
		log.Panic(err)
	}
	_api_services_ = *apiservices.apiservices
	return _api_services_
}
