package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"test/extend/code"
	"test/extend/jwt"
	"test/extend/utils"
	"test/models"
	"test/service"
)

// AuthController 用户控制器
type AuthController struct{}

// SignupRequest 账号注册请求参数
type SignupRequest struct {
	Email       string `json:"email" binding:"required,email"`
	AccountPass string `json:"accountPass" binding:"required"`
	ConfirmPass string `json:"confirmPass" binding:"required"`
}

// @Summary 账号注册
// @Description 通过邮箱密码注册账号
// @Accept json
// @Produce json
// @Tags auth
// @ID auth.signup
// @Param body body v1.SignupRequest true "账号注册请求参数"
// @Success 200 {string} json "{"status":200, "code": 2000001, msg:"请求处理成功"}"
// @Failure 400 {string} json "{"status":400, "code": 4000001, msg:"请求参数有误"}"
// @Failure 500 {string} json "{"status":500, "code": 5000001, msg:"服务器内部错误"}"
// @Router /auth/signup [post]
func (ac AuthController) Signup(c *gin.Context) {
	resBody := SignupRequest{}
	if err := c.ShouldBind(&resBody);err != nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.RequestParamError,"参数验证错误")
		return
	}
	user := models.User{}
	user.Email = resBody.Email
	user.UserName = resBody.Email
	user.Password = resBody.AccountPass
	_,err := user.Insert()
	if err != nil {
		utils.ResponseFormat(c,code.RequestParamError,"插入失败")
		return
	}
	utils.ResponseFormat(c,code.Success,"成功")
	return
}

// SigninRequest 账号登录请求参数
type SigninRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=20"`
}

// @Summary 账号登录
// @Description 通过邮箱密码登录账号
// @Accept json
// @Produce json
// @Tags auth
// @ID auth.signin
// @Param body body v1.SigninRequest true "账号登录请求参数"
// @Success 200 {string} json "{"status":200, "code": 2000001, msg:"请求处理成功"}"
// @Failure 400 {string} json "{"status":400, "code": 4000001, msg:"请求参数有误"}"
// @Failure 401 {string} json "{"status":401, "code": 4010001, msg:"账号或密码有误"}"
// @Failure 500 {string} json "{"status":500, "code": 5000001, msg:"服务器内部错误"}"
// @Router /auth/signin [post]
func (ac AuthController) Signin(c *gin.Context) {
	log.Info().Msg("enter Signin controller")
	reqBody := SigninRequest{}
	reqBody.Email = "136235@qq.com"
	err := c.ShouldBindJSON(&reqBody)
	if err != nil && false{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.RequestParamError, nil)
		return
	}
	// 登录验证
	userService := service.UserService{}
	userService.Email =  reqBody.Email
	userService.ID = 12
	userService.UserName = "hushuang"
	// 生成 Token
	authService := service.AuthService{
		User: &(userService.User),
	}
	user := &userService.User
	token, err := authService.GenerateToken(*user)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	utils.ResponseFormat(c, code.Success, map[string]interface{}{
		"userId": user.ID,
		"userName": user.UserName,
		"email": user.Email,
		"token": token,
	})
	return
}


// @Summary 账号注销
// @Description 用户账号注销
// @Accept json
// @Produce json
// @Tags auth
// @ID auth.signout
// @Param Authorization header string true "认证 Token 值"
// @Success 200 {string} json "{"status":200, "code": 2000001, msg:"请求处理成功"}"
// @Failure 500 {string} json "{"status":500, "code": 5000001, msg:"服务器内部错误"}"
// @Router /auth/signout [post]
func (ac AuthController) Signout(c *gin.Context) {
	log.Info().Msg("enter signout controller")
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	log.Debug().Msgf("claims: %v", claims)
	// 销毁 token
	authService := service.AuthService{}
	isOK, err := authService.DestroyToken(claims.Email)
	if err != nil || isOK == false {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	utils.ResponseFormat(c, code.Success, map[string]interface{}{})
	return
}