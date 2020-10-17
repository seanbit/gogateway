package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/seanbit/ginserver"
	"github.com/seanbit/gokit/foundation"
	"github.com/seanbit/gokit/validate"
	"github.com/seanbit/goserving"
	"io/ioutil"
	"net/http"
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

func TestServerBindParameter(t *testing.T) {

	LoadDatas("./test_data.json")
	var services = []serving.Registry{serving.Registry{Name: "testproduct-gateway", Rcvr: new(service), Metadata: ""}}
	serving.Serve(rpcConfig, nil, services, true)
	ginserver.Serve(httpConfig, log, func(engine *gin.Engine) {
		engine.POST("api/v1/order/pay", func(ctx *gin.Context) {
			g := ginserver.Gin{Ctx: ctx}
			var parameter = NewData("GoodsPayParameter")
			if err := g.BindParameter(parameter); err != nil {
				g.ResponseError(err)
				return
			}
			fmt.Printf("%+v\n", parameter)
			if err := validate.ValidateParameter(parameter); err != nil {
				g.ResponseError(err)
				return
			}
			if bts, err := json.Marshal(parameter); err != nil {
				g.ResponseError(err)
				return
			} else {
				g.ResponseData(string(bts))
			}
		})
	})
}

func TestPostToPay(t *testing.T)  {
	var url = fmt.Sprintf("http://localhost:%d/api/v1/order/pay", httpConfig.HttpPort)

	var user_info map[string]interface{} = make(map[string]interface{})
	user_info["user_id"] = 101
	user_info["user_name"] = "18922311056"
	user_info["email"] = "1028990481@qq.com"

	var goods1 map[string]interface{} = make(map[string]interface{})
	goods1["goods_id"] = 1001
	goods1["goods_name"] = "三只松鼠干果巧克力100g包邮"
	goods1["goods_amount"] = 1
	goods1["remark"] = ""
	var goods []interface{} = []interface{}{goods1}
	var goods_ids []int = []int{1}

	var parameter map[string]interface{} = make(map[string]interface{})
	parameter["user_info"] = user_info
	parameter["goods"] = goods
	parameter["goods_ids"] = goods_ids

	jsonStr, err := json.Marshal(parameter)
	if err != nil {
		fmt.Printf("to json error:%v\n", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		fmt.Printf("resp error:%v", err)
	} else {
		statuscode := resp.StatusCode
		hea := resp.Header
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Println(statuscode)
		fmt.Println(hea)
	}
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
