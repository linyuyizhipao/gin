package v1

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
	"test/extend/code"
	"test/extend/redis"
	"test/extend/utils"
	"test/models"
	"test/service"
	"time"
)

type AuctionController struct{}


// @Summary 竞拍时用户加价
// @Description 多用户同时对商品进行加价
// @Accept json
// @Produce json
// @Tags auction
// @ID auction.RaisePrice
// @Param body body v1.RaisePriceReq true "账号注册请求参数"
// @Success 200 {string} json "{"status":200, "code": 2000001, msg:"请求处理成功"}"
// @Failure 400 {string} json "{"status":400, "code": 4000001, msg:"请求参数有误"}"
// @Failure 500 {string} json "{"status":500, "code": 5000001, msg:"服务器内部错误"}"
// @Router /auth/signup [post]
func (a AuctionController) RaisePrice(c *gin.Context){

	log.Debug().Msg("用户抬价...")
	//接收参数并验证
	reqParam := RaisePriceReq{}
	c.ShouldBindJSON(&reqParam)
	if result,err := govalidator.ValidateStruct(reqParam);err!=nil || result != true {
		//输出错误
		log.Debug().Msg(err.Error())
		utils.ResponseFormat(c,code.RequestParamError,err.Error())
		return
	}

	//检查用户对该商品是否具备当前加价行为的资格
	userSer := service.UserService{}
	if b,err := userSer.CheckUserRaise(reqParam.Uid,reqParam.RaisePrice);err==nil || b == false{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.Success,map[string]interface{}{})
	}

	//组装用户对该商品进行该价格加价行为的消息主体转为json后发送到redis规律队列中
	if err := userSer.PushList(reqParam.Uid,reqParam.GoodsId,reqParam.RaisePrice);err!=nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.ServiceInsideError,"")
		return
	}
	//返回指定格式json给客户端，表达用户当前加价行为结果信息

	utils.ResponseFormat(c,code.Success,"竞拍成功")
	return

}
type RaisePriceReq struct {
	Uid string `json:"uid" valid:"required"`
	GoodsId string `json:"goodsId" valid:"required"`
	RaisePrice int `json:"raisePrice" valid:"required"`
}

//处理用户发送过来的加价消息，并将处理成功的消息发送到客户端
func processorRaiseMsg(){
	//找到当前待处理的所有redis队列集合
	goodsSer := service.GoodsSer

	names := goodsSer.GetAllListName()

	//遍历集合，如果该集合不为空则开启一个协程专门去处理该集合里面的抬价消息
	for _,val:= range names {
		raiseprice(val)
	}

}

//抬价消息
func raiseprice(msg string) (b bool){
	//解析json
	msgJson := [4]string{}
	goodsSer := service.GoodsSer
	json.Unmarshal([]byte(msg),&msgJson)
	goodsSer.Uid = msgJson[0]
	goodsSer.GoodsId = msgJson[1]
	goodsSer.AuctionPrice ,_= strconv.Atoi(msgJson[2])
	goodsSer.CreateTime = msgJson[3]
	//为该商品产生该用户为其加出的指定金额的价格行为记录
	if b := goodsSer.AddRaiseRecord(goodsSer.AuctionRecord);b != true {
		return
	}

	//发送订阅消息给所有在当前商品房间的所有fd(考虑用户进入到某一个商品详情，就相当于进入了一个房间，退出了详情，就相当于退出了该房间，TODO 可异步执行)
	sendMsg(goodsSer.AuctionRecord)
	return
}

//给指定房间内的所有fd发送指定消息
func sendMsg(record models.AuctionRecord){
	//找出当前商品的房间里面的所有fd


	//挨个发送消息
}





//商品10秒钟之后绑定用户行为
//
func BindGoods(){
	for{

		goodsSer := service.GoodsSer
		//获取当前参与竞拍的所有商品集合(以商品为单位数据量相对不会很大)
		raiseSet := goodsSer.GetGoodsRaiseSet()

		//循环遍历各个商品，并查看该商品的最后一次竞价记录，如果该商品最后一次竞价记录时间与当前时间对比大于或者等于10S则进行商品用户绑定行为
		for _,val := range raiseSet{
			hGoodsId := "goods:"+val
			goodsInfo,err := redis.HMGET(hGoodsId,"status","updateTime","resPrice","minPrice")
			if err != nil {
				return
			}

			if len(goodsInfo) == 2 {
				hStatus,_ := strconv.Atoi(goodsInfo[0])
				hUpdateTime,_ := strconv.ParseInt(goodsInfo[1],10,64)
				hResPrice,_:= strconv.ParseInt(goodsInfo[2],10,64)
				hMinPrice,_ := strconv.ParseInt(goodsInfo[3],10,64)

				nowTime := time.Now().Unix()
				updateTime := hUpdateTime / 1e9
				if nowTime - updateTime >= 10 && hStatus != 1{
					//lock
					if hResPrice > hMinPrice {
						//如果大于商品最低竞拍价格则改变商品状态进行最终锁定，并发送消息
						redis.HMSET(hGoodsId,"status",1)
					}else{
						//检查商品当前被用户竞拍的价格，如果小于商品最低价格则调用机器人服务去伪竞拍商品

					}

				}
			}

		}


	}

	//记上一步商品绑定(TODO 异步)
}