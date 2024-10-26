package domain

// repository中定义的表和和数据库直接映射,是用户的领域对象
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    int64
	Utime    int64
}
