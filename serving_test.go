package gateway

import (
	"context"
	"encoding/json"
	"github.com/seanbit/gokit/validate"
	serving "github.com/seanbit/goserving"
	"testing"
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



