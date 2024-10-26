package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	//添加忽略路径
	l.paths = append(l.paths, path)
	return l
}

// 校验中间件
// 验证是否有userId,没有就返回401
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		//忽略路径
		//如果是登录和注册不需要校验
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(c)
		id := sess.Get("userId")
		if id == nil {
			c.AbortWithStatus(401)
			return
		}
	}
}
