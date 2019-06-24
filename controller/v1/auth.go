package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"test/extend/code"
	"test/extend/utils"
	"test/models"
	"test/service"
)

const DEDFAULT_USER_STATUS = 0//用户默认注册时候默认的status的值
const DEFAUKT_USER_BALANCE = 0 //用户注册时候默认的账户余额

// AuthController 用户控制器
type AuthController struct{}

// SignupRequest 账号注册请求参数
type SignupRequest struct {
	Username       string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Avatar       string `json:"avatar" binding:"required"`
	Password string `json:"password" binding:"required"`
	ConfirmPass string `json:"confirmPass" binding:"required"`
}

// @Summary 账号注册
// @Description 注册自己的登录账号
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
	log.Info().Msg("开始注册了啊")
	reqBody := SignupRequest{}
	if err := c.ShouldBindJSON(&reqBody);err != nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.RequestParamError,nil)
		return
	}

	if reqBody.Password != reqBody.ConfirmPass{
		log.Error().Msg("注册确认密码二次输入不一致")
		utils.ResponseFormat(c,code.SignupPassUnmatch,nil)
		return
	}
	userModel := models.User{}
	userModel.UserName = reqBody.Username
	userModel.Password =userModel.Encryption(reqBody.Password)
	userModel.Email = reqBody.Email
	userModel.Avatar = reqBody.Avatar
	userModel.Balance = DEDFAULT_USER_STATUS
	userModel.Status = DEFAUKT_USER_BALANCE
	if _,err := userModel.Insert();err!=nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c,code.DbErrorInsert,nil)
		return
	}
	utils.ResponseFormat(c,code.Success,"")
    return
}

// SigninRequest 账号登录请求参数
type SigninRequest struct {
	Username string `json:"username" binding:"required"`
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
	if err := c.ShouldBindJSON(&reqBody);err!=nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.RequestParamError, nil)
		return
	}


	// 登录验证

	userModel := models.User{}
	userInfo,err := userModel.FindOne(map[string]interface{}{"user_name":reqBody.Username})

	if err!=nil{
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, code.DbErrorSelete, nil)
		return
	}

	reqBody.Password = userInfo.Encryption(reqBody.Password)
	if userInfo.Password != reqBody.Password{
		utils.ResponseFormat(c, code.LoginPasswordErr, nil)
		return
	}

	// 生成 Token
	authService := service.AuthService{
		User: userInfo,
	}
	token, err := authService.GenerateToken(*userInfo)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	utils.ResponseFormat(c, code.Success, map[string]interface{}{
		"userId": userInfo.ID,
		"userName": userInfo.UserName,
		"email": userInfo.Email,
		"token": token,
	})
	return
}