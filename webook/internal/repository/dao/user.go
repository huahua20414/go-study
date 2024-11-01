package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 错误
var ErrUserDulicatePhone = errors.New("手机号冲突")

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
	return dao.db.Model(&user).Where("phone = ?", user.Phone).Updates(user).Error
}

// 通过phone查询信息
func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

// 通过id查询信息
func (dao *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// 插入用户
func (dao *UserDao) Insert(ctx context.Context, u *User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.Create(&u).Error
	//判断是否有这个用户的错误
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			//邮箱冲突
			return ErrUserDulicatePhone
		}
	}
	return err
}

type User struct {
	Id       int64  `gorm:"type:bigint;primarykey;autoIncrement"`
	Phone    string `gorm:"type:varchar(20);unique"`
	Email    string `gorm:"type:varchar(70)"`
	Password string `gorm:"type:varchar(255)"`
	Ctime    int64  `gorm:"type:bigint"`
	Utime    int64  `gorm:"type:bigint"`
}
