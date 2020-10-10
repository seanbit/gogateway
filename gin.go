package gateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seanbit/gokit/encrypt"
	"github.com/seanbit/gokit/foundation"
	"github.com/seanbit/gokit/validate"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type secret_method string
const (
	//key_trace_id                      = "/seanbit/goweb/gateway/key_trace_id"
	key_ctx_trace                     = "/seanbit/goweb/gateway/key_ctx_trace"
	_                   secret_method = ""
	secret_method_rsa   secret_method = "secret_method_rsa"
	secret_method_aes   secret_method = "secret_method_aes"
	secret_method_nouse secret_method = "secret_method_nouse"
)

type CError interface {
	Code() int
	Msg() string
}

type HttpConfig struct {
	RunMode          string        `json:"-" validate:"required,oneof=debug test release"`
	WorkerId         int64         `json:"-" validate:"min=0"`
	HttpPort         int           `json:"-"`
	ReadTimeout      time.Duration `json:"read_timeout" validate:"required,gte=1"`
	WriteTimeout     time.Duration `json:"write_timeout" validate:"required,gte=1"`
	CorsAllow        bool          `json:"cors_allow"`
	CorsAllowOrigins []string      `json:"cors_allow_origins"`
	RsaOpen          bool          `json:"rsa_open"`
	RsaMap           map[string]*RsaConfig    `json:"-"`
}

/** 服务注册回调函数 **/
type GinRegisterFunc func(engine *gin.Engine)

var (
	_config   	HttpConfig
	_idWorker 	foundation.SnowId
	log *logrus.Entry
)

/**
 * 启动 api server
 * handler: 接口实现serveHttp的对象
 */
func HttpServerServe(config HttpConfig, logger logrus.FieldLogger, registerFunc GinRegisterFunc) {
	if logger == nil {
		logger = logrus.New()
	}
	log = logger.WithField("stage", "gateway")
	// config validate
	if err := validate.ValidateParameter(config); err != nil {
		log.Fatal(err)
	}
	if config.RsaOpen {
		if config.RsaMap == nil {
			log.Fatal("server http start error : secret is nil")
		}
		if err := validate.ValidateParameter(config.RsaMap); err != nil {
			log.Fatal(err)
		}
	}
	_config = config
	_idWorker, _ = foundation.NewWorker(config.WorkerId)

	// gin
	gin.SetMode(config.RunMode)
	gin.DisableConsoleColor()
	//gin.DefaultWriter = io.MultiWriter(log.Logger.Writer(), os.Stdout)

	// engine
	//engine := gin.Default()
	engine := gin.New()
	engine.Use(gin.Recovery())
	//engine.StaticFS(config.Upload.FileSavePath, http.Dir(GetUploadFilePath()))
	engine.Use(func(ctx *gin.Context) {
		var lang = ctx.GetHeader("Accept-Language")
		if  SupportLanguage(lang) == false {
			lang = LanguageZh
		}
		trace := newGinTrace(ctx)
		trace.Language = lang
		trace.TraceId = uint64(_idWorker.GetId())
		ctx.Set(key_ctx_trace, trace)
		ctx.Next()
	})
	//engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	//	// 你的自定义格式
	//	if param.ErrorMessage == "" {
	//		return fmt.Sprintf("[GIN]%s requestid:%d clientip:%s method:%s path:%s code:%d\n",
	//			param.TimeStamp.Format("2006/01/02 15:04:05"),
	//			param.Keys[key_trace_id].(uint64),
	//			param.ClientIP,
	//			param.Method,
	//			param.Path,
	//			param.StatusCode,
	//		)
	//	}
	//	return fmt.Sprintf("[GIN]%s requestid:%d clientip:%s method:%s path:%s code:%d errmsg:%s\n",
	//		param.TimeStamp.Format("2006/01/02 15:04:05"),
	//		param.Keys[key_trace_id].(uint64),
	//		param.ClientIP,
	//		param.Method,
	//		param.Path,
	//		param.StatusCode,
	//		param.ErrorMessage,
	//		)
	//}))
	if config.CorsAllow {
		if config.CorsAllowOrigins != nil {
			corscfg := cors.DefaultConfig()
			corscfg.AllowOrigins = config.CorsAllowOrigins
			engine.Use(cors.New(corscfg))
		} else {
			engine.Use(cors.Default())
		}

	}
	registerFunc(engine)
	// server
	s := http.Server{
		Addr:           fmt.Sprintf(":%d", config.HttpPort),
		Handler:        engine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(fmt.Sprintf("Listen: %v\n", err))
		}
	}()
	// signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

type Gin struct {
	Ctx *gin.Context
}

/**
 * 服务信息
 */
type ginTrace struct {
	Language 	 string
	SecretMethod secret_method `json:"secretMethod"`
	Params       []byte        `json:"params"`
	Key          []byte        `json:"key"`
	Rsa			 *RsaConfig
	TraceId 	uint64		`json:"traceId" validate:"required,gte=1"`
	UserId      uint64      `json:"userId" validate:"required,gte=1"`
	UserName    string      `json:"userName" validate:"required,gte=1"`
	UserRole	string		`json:"userRole" validate:"required,gte=1"`
}

/**
 * 请求信息创建，并绑定至context上
 */
func newGinTrace(ctx *gin.Context) *ginTrace {
	rq := &ginTrace{
		SecretMethod: secret_method_nouse,
		Params:       nil,
		Key:          nil,
		Rsa: 		  nil,
	}
	ctx.Set(key_ctx_trace, rq)
	return rq
}

/**
 * 信息获取，获取传输链上context绑定的用户请求调用信息
 */
func (g *Gin) getTrace() *ginTrace {
	obj := g.Ctx.Value(key_ctx_trace)
	if info, ok := obj.(*ginTrace); ok {
		return  info
	}
	return nil
}

/**
 * 参数绑定
 */
func (g *Gin) BindParameter(parameter interface{}) error {

	switch g.getTrace().SecretMethod {
	case secret_method_nouse:
		if err := g.Ctx.Bind(parameter); err != nil {
			return foundation.NewError(err, STATUS_CODE_INVALID_PARAMS, err.Error())
		}
		g.LogRequestParam(parameter)
		return nil
	case secret_method_aes:fallthrough
	case secret_method_rsa:
		if err := json.Unmarshal(g.getTrace().Params, parameter); err != nil {
			return foundation.NewError(err, STATUS_CODE_INVALID_PARAMS, err.Error())
		}
		g.LogRequestParam(parameter)
		return nil
	}
	return nil
}

/**
 * 响应数据，成功，原数据转json返回
 */
func (g *Gin) ResponseData(data interface{}) {
	var code = STATUS_CODE_SUCCESS
	var msg = Msg(g.getTrace().Language, code)

	switch g.getTrace().SecretMethod {
	case secret_method_nouse:
		g.LogResponseInfo(code, msg, data, "")
		g.Response(code, msg, data, "")
		return
	case secret_method_aes:
		jsonBytes, _ := json.Marshal(data)
		if secretBytes, err := encrypt.GetAes().EncryptCBC(jsonBytes, g.getTrace().Key); err == nil {
			g.LogResponseInfo(code, msg, jsonBytes, "")
			g.Response(code, msg, base64.StdEncoding.EncodeToString(secretBytes), "")
			return
		}
		g.LogResponseInfo(code, msg, data, "response data aes encrypt failed")
		g.Response(code, msg, data, "response data aes encrypt failed")
		return
	case secret_method_rsa:
		jsonBytes, _ := json.Marshal(data)
		if secretBytes, err := encrypt.GetRsa().Encrypt(g.getTrace().Rsa.ClientPubKey, jsonBytes); err == nil {
			if signBytes, err := encrypt.GetRsa().Sign(g.getTrace().Rsa.ServerPriKey, jsonBytes); err == nil {
				sign := base64.StdEncoding.EncodeToString(signBytes)
				g.LogResponseInfo(code, msg, jsonBytes, sign)
				g.Response(code, msg, base64.StdEncoding.EncodeToString(secretBytes), sign)
				return
			}
		}
		g.LogResponseInfo(code, msg, data, "response data rsa encrypt failed")
		g.Response(code, msg, data, "response data rsa encrypt failed")
		return
	}
}

/**
 * 响应数据，自定义error
 */
func (g *Gin) ResponseError(err error) {
	if e, ok := err.(foundation.Error); ok {
		msg := Msg(g.getTrace().Language, e.Code())
		g.LogResponseError(e.Code(), msg, e.Error())
		g.Response(e.Code(), msg, nil, "")
		return
	}
	g.LogResponseError(STATUS_CODE_FAILED, err.Error(), "")
	g.Response(STATUS_CODE_FAILED, err.Error(), nil, "")
}

/**
 * 响应数据
 */
func (g *Gin) Response(statusCode int, msg string, data interface{}, sign string) {
	g.Ctx.JSON(http.StatusOK, gin.H{
		"code" : statusCode,
		"msg" :  msg,
		"data" : data,
		"sign" : sign,
	})
	return
}



func (g *Gin) LogRequestParam(parameter interface{}) {
	traceId := g.getTrace().TraceId
	userId := g.getTrace().UserId
	userName := g.getTrace().UserName
	role := g.getTrace().UserRole
	apilog := log.WithFields(logrus.Fields{"traceId":traceId, "userId":userId, "userName":userName, "role":role})
	if jsonBytes, ok := parameter.([]byte); ok {
		apilog.WithField("params", string(jsonBytes)).Info("request in")
	} else if jsonBytes, err := json.Marshal(parameter); err == nil {
		apilog.WithField("params", string(jsonBytes)).Info("request in")
	} else {
		apilog.WithField("params", parameter).Info("request in")
	}
}

func (g *Gin) LogResponseInfo(code int, msg string, data interface{}, sign string) {
	traceId := g.getTrace().TraceId
	userId := g.getTrace().UserId
	userName := g.getTrace().UserName
	role := g.getTrace().UserRole
	apilog := log.WithFields(logrus.Fields{"traceId":traceId, "userId":userId, "userName":userName, "role":role, "respcode":code, "respmsg":msg, "sign":sign})

	if jsonBytes, ok := data.([]byte); ok {
		apilog.WithField("respdata", string(jsonBytes)).Info("response to")
	} else if jsonBytes, err := json.Marshal(data); err == nil {
		apilog.WithField("respdata", string(jsonBytes)).Info("response to")
	} else {
		apilog.WithField("respdata", data).Info("response to")
	}
}

func (g *Gin) LogResponseError(code int, msg string, err string) {
	traceId := g.getTrace().TraceId
	userId := g.getTrace().UserId
	userName := g.getTrace().UserName
	role := g.getTrace().UserRole
	apilog := log.WithFields(logrus.Fields{"traceId":traceId, "userId":userId, "userName":userName, "role":role, "respcode":code, "respmsg":msg})
	apilog.Info(err)
	if err != "" {
		apilog.Error(err)
	}
}

