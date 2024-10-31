package web

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go-study/day01/config"
	"go-study/day01/internal/repository"
	cache2 "go-study/day01/internal/repository/cache"
	"go-study/day01/internal/repository/dao"
	"go-study/day01/internal/service"
	"go-study/day01/internal/service/sms/email"
	"go-study/day01/internal/service/sms/tencent"
	"go-study/day01/internal/web/middleware"
	"go-study/day01/pkg/ginx/middleware/ratelimit"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func RegisterRoutes() *gin.Engine {
	//初始化数据库
	db := initDB()

	//初始化UserHandler
	u := initUser(db)

	//初始化路由,配置跨域和redis缓存判断用户是否登录并刷新缓存时间
	server := initServer()

	//注册用户路由
	u.RegisterUserRoutes(server)

	return server
}

func initServer() *gin.Engine {
	//初始化路由
	server := gin.Default()

	//配置限流中间件
	//初始化redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//跨域中间件
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost"}, // 设置允许的来源
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//允许拿什么
		AllowHeaders:  []string{"Content-Type", "Authorization", "x-jwt-token"},
		ExposeHeaders: []string{"Content-Length"},
		//是否允许携带cookie
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 缓存 CORS 设置
	}))

	//session基本使用,创建session存储在cookie中 redis存储
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
	//if err != nil {
	//	panic(err)
	//}
	store := memstore.NewStore([]byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
	server.Use(sessions.Sessions("mvsession", store))
	//检验是否有session的中间件 redis实现
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	//检验是否有session的中间件 jwt实现
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/sms").Build())
	return server
}

// 初始化要使用的userhandler对象1
func initUser(db *gorm.DB) *UserHandler {
	//初始化对象
	ud := dao.NewUserDao(db)
	//用来去redis查缓存的数据
	cache1 := cache2.NewUserCache(redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	}), time.Minute*15, time.Minute*5)
	repo := repository.NewUserRepository(ud, cache1)
	//初始化邮箱验证码配置
	m := gomail.NewMessage()
	m.SetHeader("From", "HuaHua<"+"2041436630@qq.com"+">") // 设置发件人别名
	m.SetHeader("Subject", "您的验证码")                   // 邮件主题

	PASSWORD, ok := os.LookupEnv("SMS_PASSWORD")
	if !ok {
		fmt.Println("no sms_password")
	}
	d := gomail.NewDialer(
		"smtp.qq.com",
		465,
		"2041436630@qq.com",
		PASSWORD,
	)
	//初始化腾讯云短信配置

	//创建sms对象
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		fmt.Println("SMS_SECRET_ID is empty")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		fmt.Println("SMS_SECRET_KEY is empty")
	}
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		//打印到日志中
	}
	//模板id，签名名称,创建sms接口对象
	s := tencent.NewService(c, "1400946141", "goLang科技公众号")

	el := email.NewService(d, m)
	svc := service.NewUserService(repo, el, s)
	u := NewUserHandler(svc)
	return u
}
func initDB() *gorm.DB {
	//初始化配置信息
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	//数据库同步,user表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
