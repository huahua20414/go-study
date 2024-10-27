package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-study/day01/internal/domain"
	"go-study/day01/internal/service"
	"net/http"
	"time"
)

// 错误参数
var (
	ErrUserDulicateEmail = service.ErrUserDulicateEmail
)

// 校验参数
const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc                  *service.UserService
	emailRegexExp        *regexp.Regexp
	passwordRegexPattern *regexp.Regexp
}

// 只需要传入一个service对象
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:                  svc,
		emailRegexExp:        regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexPattern: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

// 注册user路由
func (u UserHandler) RegisterUserRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("signup", u.SignUp)
	//ug.POST("login", u.Login)
	ug.POST("login", u.LoginJwt)
	ug.POST("edit", u.Edit)
	ug.GET("profile", u.ProfileJWT)

}

// 登录
func (u *UserHandler) LoginJwt(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		//c.String(200, err.Error())
		return
	}
	user, err := u.svc.Login(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		c.String(200, "用户名或者密码不对")
		return
	}
	if err != nil {
		c.String(500, "系统错误")
		return
	}
	//生成token claims对象
	claims := domain.UserClaims{
		//设置参数
		RegisteredClaims: jwt.RegisteredClaims{
			//设置1分钟的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid: user.Id,
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		//c.String(200, err.Error())
		return
	}
	user, err := u.svc.Login(c, domain.User{
		Email:    req.Email,
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
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//绑定到req对象
	if err := c.ShouldBind(&req); err != nil {
		return
	}
	isEmail, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "邮箱不正确")
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
			"密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}
	//调用service
	err = u.svc.SignUp(c, domain.User{Email: req.Email,
		Password: req.Password})
	if err == ErrUserDulicateEmail {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "邮箱冲突",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})

}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	var req EditReq
	if err := c.ShouldBind(&req); err != nil {
		return
	}
	err := u.svc.Edit(c, req.Email, req.OldPassword, req.NewPassword)
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
func (u *UserHandler) Profile(c *gin.Context) {
	c.String(200, "登录成功")
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
	println(claims.Uid)
}
