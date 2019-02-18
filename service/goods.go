package service

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"test/extend/redis"
	"test/models"
	"time"
)

const GOODS_RAISE_EXPIRE_TIME = 10
const GOODS_LUCK_STATUS = 0 //商品锁定状态，也就是被竞拍掉了的状态
var GoodsSer = goodsService{}

var listNames = []string{"goods-raise-price-list-0-30","goods-raise-price-list-30-60","goods-raise-price-list-60-100"}
// goods 对象
type goodsService struct{
	models.AuctionRecord
}
//根据商品id计算出该商品应该进入的reids的list名称
func (g *goodsService) getListNamePassGoodsId(goodsId string)(s string,err error){
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
			s = listNames[0]
			case  i > 30 && i <= 60:
			s = listNames[1]
			default :
			s = listNames[2]
	}
	return
}

//返回当前用于储存抬价的redis队列名称集合
func (g *goodsService) GetAllListName()([]string){
	return listNames
}

//用户竞拍行为入mysql,并更新goods的竞拍价格属性状态
func (g *goodsService) AddRaiseRecord(record models.AuctionRecord)(b bool){
	db := models.DB
	defer db.Close()

	hGoodsId := "goods:"+record.GoodsId

	goodsInfo,err := redis.HMGET(hGoodsId,"goodsId","status","auctionPrice","updateTime")

	if err != nil {
		return
	}

	if len(goodsInfo) != 4 {
		//商品未进入到竞拍必要的redis缓存中
		return
	}

	hStatus,_ :=strconv.Atoi(goodsInfo[1])
	hAuctionPrice := goodsInfo[2]
	hUpdateTime := goodsInfo[3]

	if hStatus != GOODS_LUCK_STATUS {
		//说明该商品已经被锁定
		return
	}

	nowTime := time.Now().Unix()
	updateTime,_ := strconv.ParseInt(hUpdateTime,10,64)
	updateTime = updateTime / 1e9
	if nowTime - updateTime > GOODS_RAISE_EXPIRE_TIME {
		//竞拍时间间隔已经超过了最大值10秒了
		return
	}

	//更新当前商品的竞拍属性状态
	auctionPrice,_ :=strconv.Atoi(hAuctionPrice)

	if auctionPrice >= record.AuctionPrice {
		//竞拍价格比商品竞拍价格要低，非法
		return
	}

	redis.HMSET(hGoodsId,"auctionPrice",hAuctionPrice,"updateTime",time.Now().UnixNano())

	//将竞拍记录储存在mysql中
	b = db.NewRecord(record)
	return
}

//获取参与竞拍行为商品的集合
func (g *goodsService)GetGoodsRaiseSet()(raiseSet []string){

	return
}