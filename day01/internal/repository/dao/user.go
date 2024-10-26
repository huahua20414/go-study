package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 错误
var ErrUserDulicateEmail = errors.New("邮箱冲突")

type UserDao struct {
	db *gorm.DB
}

// 用来初始化userdao实例
func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

// 更新用户信息
func (dao *UserDao) Updates(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Utime = now
	return dao.db.Model(&user).Where("email = ?", user.Email).Updates(user).Error
}

// 通过email查询信息
func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

// 插入用户
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.Create(&u).Error
	//判断是否有这个用户的错误
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			//邮箱冲突
			return ErrUserDulicateEmail
		}
	}
	return err
}

type User struct {
	Id       int64  `gorm:"type:bigint;primarykey;autoIncrement"`
	Email    string `gorm:"type:varchar(70);unique"`
	Password string `gorm:"type:varchar(255)"`
	Ctime    int64  `gorm:"type:bigint"`
	Utime    int64  `gorm:"type:bigint"`
}
