package service

import (
	"github.com/gin-gonic/gin/json"
	"github.com/rs/zerolog/log"
	"test/extend/redis"
	"test/models"
	"time"
)

// UserService 用户服务层逻辑
type UserService struct{
	models.User
}

//检查用户是否具备当前加价行为资格
func (us *UserService) CheckUserRaise(uid string,raisePrice int) (b bool, err error) {

	condition:=map[string]interface{}{"uid":uid}
	userInfo,err := us.FindOne(condition)
	if err != nil{
		log.Error().Msg(err.Error())
		return
	}
	b = userInfo.Balance > raisePrice
	return
}
//用户加价行为入redis队列
func (us *UserService) PushList(uid string,goodsId string ,raise int)(err error){
	listKey,err := GoodsSer.getListNamePassGoodsId(goodsId)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	createTime := time.Now().UnixNano()
	data :=[4]string{uid,goodsId,string(raise),string(createTime)}
	j,err :=json.Marshal(data)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	err = redis.SADD(listKey,string(j))
	return
}