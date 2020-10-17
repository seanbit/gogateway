package gateway

import (
	"fmt"
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
	LoadDatas("./test_data.json")
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