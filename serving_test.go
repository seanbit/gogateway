package gateway

import (
	"context"
	"encoding/json"
	"github.com/seanbit/gokit/foundation"
	"github.com/seanbit/gokit/validate"
	"github.com/seanbit/goserving"
	"testing"
	"time"
)

func TestServing(t *testing.T) {
	// serving
	var services = []serving.Registry{
		serving.Registry{
			Name:     "test-order",
			Rcvr:     orderService{},
			Metadata: "",
		},
	}
	serving.Serve(rpcConfig, log, services, false)
}

type orderService struct {
}
var OrderService = new(orderService)

func (this *orderService) GoodsPay(ctx context.Context, parameter *GoodsPayParameter, resp *string) error {
	if err := validate.ValidateParameter(parameter); err != nil {
		return err
	}
	if bts, err := json.Marshal(parameter); err != nil {
		return err
	} else {
		*resp = string(bts)
	}
	return nil
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



