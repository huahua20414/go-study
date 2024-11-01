//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"go-study/wire/repository"
	"go-study/wire/repository/dao"
)

func InitReposity() *repository.UserRepository {
	//这个方法里面传入各种初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
