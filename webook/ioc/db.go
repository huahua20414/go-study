package ioc

import (
	"go-study/webook/config"
	"go-study/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
