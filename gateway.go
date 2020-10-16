package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/seanbit/ginserver"
	"github.com/seanbit/goserving"
	"github.com/sirupsen/logrus"
)

var log logrus.FieldLogger
var middlewares = make(map[string]gin.HandlerFunc)

func Serve(pname string, rpcConfig serving.RpcConfig, httpConfig ginserver.HttpConfig, logger logrus.FieldLogger, dataPath, apiPath string)  {
	if logger == nil {
		logger = logrus.New()
	}
	log = logger
	LoadDatas(dataPath)
	apiServices := LoadApis(apiPath)
	var services = []serving.Registry{serving.Registry{Name: pname + "-gateway", Rcvr: new(service), Metadata: ""}}
	serving.Serve(rpcConfig, nil, services, true)
	ginserver.Serve(httpConfig, log, func(engine *gin.Engine) {
		for _, service := range apiServices {
			apiRegister(engine, nil, service)
		}
	})
}

func apiRegister(engine *gin.Engine, group *gin.RouterGroup, service APIService)  {
	// handlers
	var handlers = make([]gin.HandlerFunc, 0)
	if len(service.MiddelWares) > 0 {
		for _, mwName := range service.MiddelWares {
			if handler, ok := middlewares[mwName]; ok {
				handlers = append(handlers, handler)
			}
		}
	}
	// register group
	isGroup := len(service.Services) > 0
	if isGroup {
		var parentGroup *gin.RouterGroup
		if group == nil {
			parentGroup = engine.Group(service.Path, handlers...)
		} else {
			parentGroup = group.Group(service.Path, handlers...)
		}
		for _, subService := range service.Services {
			apiRegister(engine, parentGroup, subService)
		}
		return
	}
	// register api
	var apiHandler = func(c *gin.Context) {
		g := ginserver.Gin{Ctx: c}
		parameter := NewData(service.Do.RpcRequest)
		response := NewData(service.Do.RpcResponse)
		serving.TraceBind(g.Ctx, g.Trace().TraceId, g.Trace().UserId, g.Trace().UserName, g.Trace().UserRole)
		err := serving.Call(service.Do.RpcService, service.Do.RpcServer, g.Ctx, service.Do.RpcMethod, parameter, response)
		if err != nil {
			g.ResponseError(err)
			return
		}
		g.ResponseData(response)
	}
	handlers = append(handlers, apiHandler)
	switch service.Do.HttpMethod {
	case HTTPMethodGET:
		group.GET(service.Path, handlers...)
	case HTTPMethodPOST:
		group.POST(service.Path, handlers...)
	case HTTPMethodPUT:
		group.PUT(service.Path, handlers...)
	case HTTPMethodDELETE:
		group.DELETE(service.Path, handlers...)
	}
}

func MiddleWareSet(mw map[string]gin.HandlerFunc)  {
	middlewares = mw
}

func MiddleWareAdd(name string, handler gin.HandlerFunc) {
	middlewares[name] = handler
}



type service struct {
}
func (this *service) Ping(ctx context.Context, ping *string, pong *string) error {
	*pong = *ping
	return nil
}