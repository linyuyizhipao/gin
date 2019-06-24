package service

import (
	"github.com/rs/zerolog/log"
	"test/models"
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

