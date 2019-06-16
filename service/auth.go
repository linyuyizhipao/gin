package service

import (
	"time"

	"test/extend/conf"
	"test/extend/jwt"
	"test/extend/redis"
	"test/models"
	goJWT "github.com/dgrijalva/jwt-go"
)

const LOGIN_REDIS_PREFIX_KEY_NAME = "TOKEN:"
const LOGIN_REDIS_KEY_VALID = 60 * 60
// AuthService 认证相关
type AuthService struct {
	User *models.User
}

// GenerateToken 生成 Token
func (as *AuthService) GenerateToken(user models.User) (string, error) {
	jwtInstance := jwt.NewJWT()
	nowTime := time.Now()
	expireTime := time.Duration(conf.ServerConf.JWTExpire)
	claims := jwt.CustomClaims{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		StandardClaims: goJWT.StandardClaims{
			ExpiresAt: nowTime.Add(expireTime * time.Hour).Unix(),
			Issuer:    "hugo",
		},
	}
	// 创建token
	token, err := jwtInstance.CreateToken(claims)
	if err != nil {
		return "", err
	}

	// 设置redis缓存
	const hourSecs int = LOGIN_REDIS_KEY_VALID
	redis.Set(LOGIN_REDIS_PREFIX_KEY_NAME+user.Email, token, conf.ServerConf.JWTExpire * hourSecs)
	return token, nil
}

// DestroyToken 销毁 Token
func (as *AuthService) DestroyToken(email string) (bool, error) {
	return redis.Del(LOGIN_REDIS_PREFIX_KEY_NAME+email)
}
