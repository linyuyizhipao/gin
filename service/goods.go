package service

import (
	"github.com/rs/zerolog/log"
	"strconv"
)

var GoodsSer = GoodsService{}

// goods 对象
type GoodsService struct{
}
//根据商品id计算出该商品应该进入的reids的list名称
func (g *GoodsService) getListNamePassGoodsId(goodsId string)(s string,err error){
	if len(goodsId) != 20 {
		log.Error().Msg("商品id长度不正确")
		return
	}
	s = goodsId[0:1]
	i,err :=strconv.ParseInt(s,10,64)
	if err != nil{
		log.Error().Msg(err.Error())
		return
	}
	switch  {
		case  i > 0 && i <= 30:
			s = "goods-raise-price-list-0-30"
			case  i > 30 && i <= 60:
			s = "goods-raise-price-list-30-60"
			default :
			s = "goods-raise-price-list-60-100"
	}
	return
}