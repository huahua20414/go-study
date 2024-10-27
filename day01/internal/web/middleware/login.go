package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
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
		//检验userId是否有内容如果是空就返回401并且结束,否则就让他进去
		sess := sessions.Default(c)
		id := sess.Get("userId")
		if id == nil {
			c.AbortWithStatus(401)
			return
		}
		//进去,刷新用户时间间隔
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli()
		//说明还没有刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		//updateTime是有的,断言
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			c.AbortWithStatus(401)
			return
		}
		//如果大于一分钟就让他刷新
		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

	}
}
