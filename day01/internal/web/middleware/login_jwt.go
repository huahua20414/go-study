package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-study/day01/internal/domain"
	"log"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	//添加忽略路径
	l.paths = append(l.paths, path)
	return l
}

// 校验中间件
// 验证是否有userId,没有就返回401
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		//忽略路径
		//如果是登录和注册不需要校验
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		//我现在用jwt来检验
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			//说明没有登录
			c.AbortWithStatus(401)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			c.AbortWithStatus(401)
			return
		}
		//获取token
		var claims = &domain.UserClaims{}
		//获取除了Bearer那一坨字符串
		tokenStr := segs[1]
		//解析token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			//密钥
			return []byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"), nil
		})
		if claims.UserAgent != c.Request.UserAgent() {
			//严重的安全问题
			c.AbortWithStatus(401)
			return
		}
		//有错误 过期了 uid为初始值0
		if err != nil || !token.Valid || claims.Uid == 0 {
			//解析失败 没有登陆
			c.AbortWithStatus(401)
			return
		}
		//如果没有过期
		//每十秒钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			//claims是用来生成token的有关对象,设置一个新的token一分钟有效期
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err := token.SignedString([]byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
			if err != nil {
				log.Println("续约失败")
				return
			}
			//重新设置一个token
			c.Header("x-jwt-token", tokenStr)
		}
		//err为空即解析成功用户可以登录,每次请求接口时都将Authorization解析后并设置claims,里面有id等
		c.Set("claims", claims)
	}
}
