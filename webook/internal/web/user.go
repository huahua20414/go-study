package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go-study/webook/internal/domain"
	"go-study/webook/internal/service"
	"net/http"
	"time"
)

// 错误参数
var (
	ErrUserDulicatePhone = service.ErrUserDulicatePhone
	ErrBusyVerification  = service.ErrBusyVerification
	ErrVerificationFalse = service.ErrVerificationFalse
	ErrUserHave          = service.ErrUserHave
	ErrUserNil           = service.ErrUserNil
)

// 校验参数
const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	phoneRegexPattern    = `^1[3-9]\d{9}$`
)

type UserHandler struct {
	svc                  *service.UserService
	emailRegexExp        *regexp.Regexp
	passwordRegexPattern *regexp.Regexp
	phoneRegexExp        *regexp.Regexp
}

// 只需要传入一个service对象
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:                  svc,
		emailRegexExp:        regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexPattern: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRegexExp:        regexp.MustCompile(phoneRegexPattern, regexp.None),
	}
}

// 注册user路由
func (u UserHandler) RegisterUserRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("signup", u.SignUp)
	//ug.POST("login", u.Login)
	ug.POST("login", u.LoginJwt)
	ug.POST("/:action/code/send", u.Sms)
	ug.POST("edit", u.Edit)
	ug.GET("profile", u.ProfileJWT)
}

func (u *UserHandler) LoginSms(c *gin.Context)  {}
func (u *UserHandler) ForgetSms(c *gin.Context) {}

// 验证码发送
func (u *UserHandler) Sms(c *gin.Context) {

	type SmsReq struct {
		Phone string `json:"phone"`
	}
	var req SmsReq
	if err := c.ShouldBind(&req); err != nil {
		c.String(500, "系统错误")
		return
	}
	isPhone, err := u.phoneRegexExp.MatchString(req.Phone)
	if err != nil {
		c.String(500, "系统错误")
		return
	}
	if !isPhone {
		c.String(http.StatusOK, "手机号格式不正确")
		return
	}

	//这里判断codeType
	u1 := domain.User{Phone: req.Phone}
	codeType := c.Params[0].Value
	if codeType == "register_sms" {
		u1.CodeType = "register"
		err = u.svc.RegisterSms(c, u1)
	} else if codeType == "login_sms" {
		u1.CodeType = "login"
		//todo:这里注释掉了并且login_sms service还没有写
		//err = u.svc.LoginSms(c, u1)
	} else if codeType == "forget_sms" {
		u1.CodeType = "forget"
		err = u.svc.ForgetSms(c, u1)
	}

	if err == ErrBusyVerification {
		c.String(200, "一分钟后再试")
		return
	}
	if err == ErrUserHave {
		c.String(200, "此用户已经注册")
		return
	}
	if err == ErrUserNil {
		c.String(200, "没有此用户")
		return
	}
	if err != nil {
		c.String(500, "系统错误")
		return
	}
	c.String(200, "发送验证码成功")

}

// 登录
func (u *UserHandler) LoginJwt(c *gin.Context) {
	type LoginReq struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		//c.String(200, err.Error())
		return
	}
	user, err := u.svc.Login(c, domain.User{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		c.String(500, "用户名或者密码不对")
		return
	}
	//生成token claims对象
	claims := domain.UserClaims{
		//设置参数
		RegisteredClaims: jwt.RegisteredClaims{
			//设置60分钟的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Uid:       user.Id,
		UserAgent: c.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密
	tokenStr, err := token.SignedString([]byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
	if err != nil {
		c.String(500, "系统错误")
		return
	}
	c.Header("x-jwt-token", tokenStr)
	c.String(200, "登录成功")
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		//c.String(200, err.Error())
		return
	}
	user, err := u.svc.Login(c, domain.User{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		c.String(200, "用户名或者密码不对")
		return
	}
	if err != nil {
		c.String(200, "系统错误")
		return
	}
	//设置cookie
	sess := sessions.Default(c)
	//不是唯一的key不推荐
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	//cookie保存后会生成一个mvsession(自己起的名)的cookie
	if err := sess.Save(); err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, "登录成功")
}

// 注册
func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Phone           string `json:"phone"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
		Verification    string `json:"verification"`
	}
	var req SignUpReq
	//绑定到req对象
	if err := c.ShouldBind(&req); err != nil {
		return
	}
	isPhone, err := u.phoneRegexExp.MatchString(req.Phone)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPhone {
		c.String(http.StatusOK, "手机号不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入的密码不相同")
		return
	}

	isPassword, err := u.passwordRegexPattern.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		c.String(http.StatusOK,
			"密码必须包含数字、字母，并且长度不能小于 8 位")
		return
	}
	//调用service
	err = u.svc.SignUp(c, domain.User{Phone: req.Phone,
		Password: req.Password, Verification: req.Verification, CodeType: "register"})
	if err == ErrUserDulicatePhone {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "手机号已被注册",
		})
		return
	} else if err == ErrVerificationFalse {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "验证码错误",
		})
		return
	} else if err == redis.Nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "验证码错误",
		})
		return
	}
	if err != nil {
		c.String(500, "系统错误")
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})

}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Phone       string `json:"phone"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	var req EditReq
	if err := c.ShouldBind(&req); err != nil {
		return
	}
	err := u.svc.Edit(c, req.Phone, req.OldPassword, req.NewPassword)
	if err == service.ErrInvalidPassword {
		c.String(200, "密码错误")
		return
	}
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, "修改成功")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		//很奇怪的错误,因为登录的时候已经设置了
		ctx.String(200, "系统错误")
		return
	}
	claims, ok := c.(*domain.UserClaims)
	if !ok {
		ctx.String(200, "系统错误")
		return
	}
	//在这里调用service接口去查个人信息
	profile, err := u.svc.Profile(ctx, claims.Uid)
	if err != nil {
		//日志
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"profile": profile,
	})
}
