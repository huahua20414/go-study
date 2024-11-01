package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-study/webook/internal/web"
	"go-study/webook/internal/web/middleware"
	"go-study/webook/pkg/ginx/middleware/ratelimit"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterUserRoutes(server)
	return server
}
func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:3000"}, // 设置允许的来源
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			//允许拿什么
			AllowHeaders:  []string{"Content-Type", "Authorization", "x-jwt-token"},
			ExposeHeaders: []string{"Content-Length"},
			//是否允许携带cookie
			AllowCredentials: true,
			MaxAge:           12 * time.Hour, // 缓存 CORS 设置
		}),
		sessions.Sessions("mvsession",
			memstore.NewStore([]byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/register_sms/code/send").
			IgnorePaths("/users/forget_sms/code/send").Build(),
	}
}
