package gateway

import (
	"github.com/natefinch/lumberjack"
	"github.com/seanbit/ginserver"
	"github.com/seanbit/gokit/foundation"
	"github.com/seanbit/goserving"
	"runtime"
	"testing"
	"time"
)

const ServerName = "gatewaytest"

func TestGateway(t *testing.T) {
	rpcConfig.RpcPort = 6185
	// concurrent
	runtime.GOMAXPROCS(runtime.NumCPU())
	// gateway
	Serve(ServerName, rpcConfig, httpConfig, nil, "./data.json", "./api.json")
}

func TestPost(t *testing.T) {

}

var rpcConfig = serving.RpcConfig{
	RunMode:              foundation.RUN_MODE_DEBUG,
	RpcPort:              9001,
	RpcPerSecondConnIdle: 500,
	ReadTimeout:          30 * time.Second,
	WriteTimeout:         30 * time.Second,
	TokenAuth:            false,
	Token:                nil,
	TlsAuth:              false,
	Tls:                  nil,
	WhitelistAuth:        false,
	Whitelist:            nil,
	Registry:             &serving.EtcdRegistry{
		EtcdRpcUserName: "root",
		EtcdRpcPassword: "etcd.user.root.pwd",
		EtcdRpcBasePath: "seanbit/serving/rpc/credence",
		EtcdEndPoints:   []string{"127.0.0.1:2379"},
	},
}

var httpConfig = ginserver.HttpConfig{
	RunMode:          foundation.RUN_MODE_DEBUG,
	WorkerId:         0,
	HttpPort:         6088,
	ReadTimeout:      30 * time.Second,
	WriteTimeout:     30 * time.Second,
	CorsAllow:        false,
	CorsAllowOrigins: nil,
	RsaOpen:          false,
	RsaMap:           nil,
}

var info_wirter = &lumberjack.Logger{
	Filename:   "./test_log_serving.log",
	MaxSize:    100,
	MaxBackups: 10,
	MaxAge:     30,
	Compress:   false,
}
var warn_wirter = &lumberjack.Logger{
	Filename:   "./test_log_serving.log",
	MaxSize:    100,
	MaxBackups: 10,
	MaxAge:     30,
	Compress:   false,
}
var error_wirter = &lumberjack.Logger{
	Filename:   "./test_log_error.log",
	MaxSize:    100,
	MaxBackups: 10,
	MaxAge:     30,
	Compress:   false,
}
