package domain

import "github.com/golang-jwt/jwt/v5"

// repository中定义的表和和数据库直接映射,是用户的领域对象
// 这是用户表
type User struct {
	Id           int64
	Email        string
	Phone        string
	Password     string
	Ctime        int64
	Utime        int64
	Verification string
	CodeType     string `json:"codeType"`
}

// 存到token里的东西,和用户认证有关
type UserClaims struct {
	jwt.RegisteredClaims
	//声明你自己要放在token里面的数据
	Uid       int64
	UserAgent string
}
