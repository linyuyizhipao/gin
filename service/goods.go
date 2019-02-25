package service

import (
	"errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"test/extend/redis"
	"test/models"
	"time"
)

const GOODS_RAISE_EXPIRE_TIME = 10
const GOODS_LUCK_STATUS = 0 //商品锁定状态，也就是被竞拍掉了的状态
const GOODS_AUCTION_IDS = "goods-auction-ids"


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
func (g *goodsService) AddRaiseRecord(record models.AuctionRecord)(err error){
	db := models.DB
	defer db.Close()

	hGoodsId := "goods:"+record.GoodsId

	goodsInfo,err := redis.HMGET(hGoodsId,"goodsId","status","auctionPrice","updateTime")

	if err != nil {
		err = errors.New("商品抬价信息在redis中未找到")
		return
	}

	if len(goodsInfo) != 4 {
		//商品未进入到竞拍必要的redis缓存中
		err = errors.New("商品redis找到的抬价信息不为默认规则的4个元素的数组")
		return
	}

	hStatus,_ :=strconv.Atoi(goodsInfo[1])
	hAuctionPrice := goodsInfo[2]
	hUpdateTime := goodsInfo[3]

	//检查该商品状态是否满足竞拍状态
	if hStatus != GOODS_LUCK_STATUS {
		//说明该商品已经被锁定
		err = errors.New("商品已经被锁定了，不能参与竞拍")
		return
	}

	//检查当前商品时间上是否还满足竞拍要求
	nowTime := time.Now().Unix()
	updateTime,_ := strconv.ParseInt(hUpdateTime,10,64)
	updateTime = updateTime / 1e9
	if nowTime - updateTime > GOODS_RAISE_EXPIRE_TIME {
		//竞拍时间间隔已经超过了最大值10秒了
		err = errors.New("商品竞拍时间间隔已经大于10秒了，虽然未被及时锁定，但是不再参与竞拍")
		return
	}

	//以上都满足了之后才能，进行以下步骤。1.更新商品redis上的竞拍属性，添加mysql的竞拍记录

	redis.HMSET(hGoodsId,"auctionPrice",hAuctionPrice,"updateTime",time.Now().UnixNano())

	//将竞拍记录储存在mysql中
	db.NewRecord(record)
	return
}

//获取参与竞拍行为商品的集合
//goods-auction-ids
func (g *goodsService)GetGoodsRaiseSets()(raiseSet []string,err error){

	return redis.SMEMBERS(GOODS_AUCTION_IDS)
}