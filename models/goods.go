package models

//商品基本属性信息表
type Goods struct {
	GoodsId string
	Names string
	Desc string
	Type int //商品类型
	Price int //商品指导价格
}

//商品竞拍有关必要属性
//某个商品在什么时刻被谁以什么价格竞拍到了
type AuctionGoods struct {
	GoodsId string
	BeginAuctionTime string //商品竞拍开始时间
	minPrice int //竞拍允许的最小价格
	Status int8 //0：未锁定可被继续竞拍叫价；1：已被用户锁定
	Uid string //用户uid
	ResPrice int //最终的竞拍价格
	LockTime string //锁定时间
}

//商品竞拍记录信息结构体
//某个商品在什么时候被谁以什么价格竞拍了
type AuctionRecord struct {
    GoodsId string
    Uid string
    AuctionPrice int
    CreateTime string //产生时间
}