package gateway

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	UserId int						`json:"user_id" validate:"required,min=1"`
	UserName string					`json:"user_name" validate:"required,eorp"`
	Email string					`json:"email" validate:"required,email"`
}

type GoodsPayInfoParameter struct {
	GoodsId int			`json:"goods_id" validate:"required,min=1"`
	GoodsName string	`json:"goods_name" validate:"required,gte=1"`
	GoodsAmount int		`json:"goods_amount" validate:"required,min=1"`
	Remark string 		`json:"remark" validate:"gte=0"`
}

type GoodsPayParameter struct {
	UserInfo *User					`json:"user_info" validate:"required"`
	Goods []*GoodsPayInfoParameter	`json:"goods" validate:"required,gte=1,dive,required"`
	GoodsIds []int				`json:"goods_ids" validate:"required,gte=1,dive,min=1"`
}

func TestDataNew(t *testing.T) {
	DataDefines("./test_data.json")
	parameter := NewData("GoodsPayParameter")
	fmt.Printf("%+v\n", parameter)
}

func TestParse(t *testing.T) {
	var someDt = "{NewAuthParameter}"
	params := ftRegexp.FindStringSubmatch(someDt)
	for idx, p := range params {
		fmt.Printf("%d---%s\n", idx, p)
	}
	var someDts = "[Goods]"
	params = ftsRegexp.FindStringSubmatch(someDts)
	for idx, p := range params {
		fmt.Printf("%d---%s\n", idx, p)
	}
}

func TestNewDataParse(t *testing.T) {
	var dtname = "string"
	if ft, ok := parseBaseType(dtname); ok {
		var val = reflect.New(ft)
		fmt.Printf("%+v", val)
		//return reflect.Zero(ft)
	} else if params := ftRegexp.FindStringSubmatch(dtname); len(params) == 2 && params[0] == dtname {
		typeName := params[1]
		if dt, ok := _data_types_[typeName]; ok {
			ft = dt
			var val = reflect.New(ft)
			fmt.Printf("%+v", val)
		} else {
			log.Fatalf("data parse type err: could not found typename {%s} before this data {%s}", typeName, "somedata")
		}
	} else if params = ftsRegexp.FindStringSubmatch(dtname); len(params) == 2 && params[0] == dtname { // slice shoud handle
		typeName := params[1]
		if dt, ok := parseBaseType(typeName); ok {
			ft = reflect.SliceOf(dt)
			var val = reflect.MakeSlice(dt, 0, 0)
			fmt.Printf("%+v", val)
		} else if dt, ok := _data_types_[typeName]; ok {
			ft = reflect.SliceOf(dt)
			var val = reflect.MakeSlice(dt, 0, 0)
			fmt.Printf("%+v", val)
		} else {
			log.Fatalf("data parse type err: could not found typename [%s] before this data {%s}", typeName, dtname)
		}
	} else {
		log.Fatalf("data parse type err: could not parse type with %s in this data %s", dtname, "somedata")
	}
}