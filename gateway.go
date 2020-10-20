package gateway

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/seanbit/ginserver"
	"github.com/seanbit/goserving"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	log logrus.FieldLogger
	middlewares = make(map[string]gin.HandlerFunc)
	definedHandlers = make(map[string]gin.HandlerFunc)
)

func MiddleWareSet(mw map[string]gin.HandlerFunc)  {
	middlewares = mw
}

func MiddleWareAdd(name string, handler gin.HandlerFunc) {
	middlewares[name] = handler
}

func DefinedHandlersSet(ddHandlers map[string]gin.HandlerFunc)  {
	definedHandlers = ddHandlers
}

func DefinedHandlerAdd(path string, handler gin.HandlerFunc) {
	definedHandlers[path] = handler
}

func Serve(pname string, rpcConfig serving.RpcConfig, httpConfig ginserver.HttpConfig, logger logrus.FieldLogger, dataPath, apiPath string)  {
	if logger == nil {
		logger = logrus.New()
	}
	log = logger
	DataDefines(dataPath)
	apiServices := ApiDefines(apiPath)
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
	var apiHandler gin.HandlerFunc
	if group != nil {
		path := group.BasePath() + "/" + service.Path
		path = strings.ReplaceAll(path, "//", "/")
		if handler, ok := definedHandlers[path]; ok {
			apiHandler = handler
			goto HANDLER_REGISTER
		}
	}
	apiHandler = func(c *gin.Context) {
		g := ginserver.Gin{Ctx: c}
		parameter := NewData(service.Do.RpcParameter)
		if parameter == nil {
			g.ResponseError(errors.New("api data parameter init nil in path:" + service.Path + " dataType:" + service.Do.RpcParameter))
			return
		}
		if err := g.BindParameter(parameter); err != nil {
			g.ResponseError(err)
			return
		}
		response := NewData(service.Do.RpcResponse)
		if response == nil {
			g.ResponseError(errors.New("api data response init nil in path:" + service.Path + " dataType:" + service.Do.RpcResponse))
			return
		}
		serving.TraceBind(g.Ctx, g.Trace().TraceId, g.Trace().UserId, g.Trace().UserName, g.Trace().UserRole)
		err := serving.Call(service.Do.RpcService, service.Do.RpcServer, g.Ctx, service.Do.RpcMethod, parameter, response)
		if err != nil {
			g.ResponseError(err)
			return
		}
		g.ResponseData(response)
	}
	HANDLER_REGISTER:
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



type service struct {
}
func (this *service) Ping(ctx context.Context, ping *string, pong *string) error {
	*pong = *ping
	return nil
}