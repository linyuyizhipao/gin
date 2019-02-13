package middleware

import (
	"github.com/gin-gonic/gin"
	"test/extend/utils"
	"test/extend/code"
	"test/extend/jwt"
	"test/extend/redis"
	"github.com/rs/zerolog/log"
)

// JWTAuth Token 认证中间件
// 这个jwt的正常运转，需要我们在填充 jwt.CustomClaims结构体的各字段，当然它最好或者说必须得跟用户数据库user表进行映射关联吧
// 这个jwt它还会去redis去根据cliam中的email值去取对应的token用于对比当前token，当然你也可以通过其他方式
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取 Authorization token 值
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// 获取不到 Authorization 报：请求未携带Token,无权访问
			utils.ResponseFormat(c, code.TokenNotFound, nil)
			c.Abort()
			return
		}

		// 获取到jwt实例 用于关于Token方面的各种操作
		jwtInstance := jwt.NewJWT()
		//解析token
		claims, err := jwtInstance.ParseToken(token)
		if err != nil {
			// 未能正常解析 Token，则报：token认证失败
			utils.ResponseFormat(c, code.TokenInvalid, nil)
			c.Abort()
			return
		}

		// TODO 获取缓存中的Token信息
		tokenCache, err := redis.Get("TOKEN:"+claims.Email)
		if err != nil {
			log.Error().Msgf("jwt auth redis get: %v", err.Error())
			utils.ResponseFormat(c, code.ServiceInsideError, nil)
			c.Abort()
			return
		}

		// 用户注销或token失效
		if tokenCache != token {
			log.Error().Msg("user signout or token invalid")
			utils.ResponseFormat(c, code.TokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
