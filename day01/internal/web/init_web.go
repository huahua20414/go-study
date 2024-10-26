package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go-study/day01/internal/repository"
	"go-study/day01/internal/repository/dao"
	"go-study/day01/internal/service"
	"go-study/day01/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func RegisterRoutes() *gin.Engine {
	//初始化数据库
	db := initDB()

	//初始化UserHandler
	u := initUser(db)

	//初始化路由,配置跨域
	server := initServer()

	//注册用户路由
	u.RegisterUserRoutes(server)

	return server
}

func initServer() *gin.Engine {
	//初始化路由
	server := gin.Default()
	//跨域中间件
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 设置允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 缓存 CORS 设置
	}))

	//session基本使用,创建session存储在cookie中，密钥是secret
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mvsession", store))
	//检验是否有session的中间件
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return server
}

// 初始化要使用的userhandler对象
func initUser(db *gorm.DB) *UserHandler {
	//初始化对象
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := NewUserHandler(svc)
	return u
}
func initDB() *gorm.DB {
	//初始化配置信息
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/gorm"))
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
