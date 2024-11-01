package service

import (
	"context"
	"errors"
	"go-study/webook/internal/domain"
	"go-study/webook/internal/repository"
	"go-study/webook/internal/service/sms/email"
	"go-study/webook/internal/service/sms/tencent"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// 错误
var (
	//邮箱冲突
	ErrUserDulicatePhone     = repository.ErrUserDulicatePhone
	ErrInvalidUserOrPassword = errors.New("账号或密码不对")
	ErrInvalidPassword       = errors.New("密码不对")
	ErrUserHave              = errors.New("此用户已经注册")
	ErrUserNil               = errors.New("没有此用户")
	ErrBusyVerification      = errors.New("一分钟后再试")
	ErrVerificationFalse     = errors.New("验证码错误")
)

type UserServiceInterface interface {
	ForgetSms(ctx context.Context, u domain.User) error
	RegisterSms(ctx context.Context, u domain.User) error
	Edit(ctx context.Context, phone string, op string, np string) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type UserService struct {
	repo    repository.UserRepository
	client  *email.Service
	tclient *tencent.Service
}

// 初始化
func NewUserService(repo repository.UserRepository, client *email.Service, tc *tencent.Service) UserServiceInterface {
	return &UserService{repo: repo,
		client: client, tclient: tc}
}
func (svc *UserService) ForgetSms(ctx context.Context, u domain.User) error {
	//先去查是否有这个用户如果没有提示没有此用户
	_, err := svc.repo.FindByPhone(ctx, u.Phone)
	if err != nil {
		return ErrUserNil
	}
	//发送验证码
	//查一下是否在一分钟内发过
	user, err := svc.repo.GetVerification(ctx, u)
	if err == nil {
		//有缓存说明已经请求过一次了 检验utime-now是否大于一分钟
		errand := time.Now().Unix() - user.Utime
		if errand < 60 {
			//小于一分钟
			return ErrBusyVerification
		}
	}
	if err != nil {
		//没缓存向下进行
		//todo:打日志
	}
	rand.Seed(time.Now().UnixNano())     // 以当前时间为随机数种子
	number := rand.Intn(900000) + 100000 // 生成100000到999999之间的随机数
	code := strconv.Itoa(number)         // 生成六位数验证码
	newNumber := "+86" + u.Phone         //拼接发送的电话号
	//发送电话号
	Cases := struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		//模板id
		tplId: "2297671",
		//验证码
		params: []string{code},
		// 你的手机号码
		numbers: []string{newNumber},
	}

	if err := svc.tclient.Send(ctx, Cases.tplId, Cases.params, Cases.numbers...); err != nil {
		//打个日志
		return err
	}
	u.Verification = code
	//验证码放入缓存
	err = svc.repo.SetVerification(ctx, u)
	return err
}

func (svc *UserService) RegisterSms(ctx context.Context, u domain.User) error {
	//先去查是否有这个用户如果有返回有这个用户的错误
	_, err := svc.repo.FindByPhone(ctx, u.Phone)
	//有这个用户
	if err == nil {
		return ErrUserHave
	}
	if err != nil {
		//todo:这里打个日志
	}
	//没有,发送验证码
	//查一下是否在一分钟内发过
	user, err := svc.repo.GetVerification(ctx, u)
	if err == nil {
		//有缓存说明已经请求过一次了 检验utime-now是否大于一分钟
		//todo:这里要加时间验证我注释掉先
		errand := time.Now().Unix() - user.Utime
		if errand < 60 {
			//小于一分钟
			return ErrBusyVerification
		}
	}
	if err != nil {
		//没缓存向下进行
		//todo:打日志
	}
	rand.Seed(time.Now().UnixNano())     // 以当前时间为随机数种子
	number := rand.Intn(900000) + 100000 // 生成100000到999999之间的随机数
	code := strconv.Itoa(number)         // 生成六位数验证码
	//newNumber := "+86" + u.Phone         //拼接发送的电话号
	////发送电话号
	//Cases := struct {
	//	name    string
	//	tplId   string
	//	params  []string
	//	numbers []string
	//	wantErr error
	//}{
	//	//模板id
	//	tplId: "2297671",
	//	//验证码
	//	params: []string{code},
	//	// 你的手机号码
	//	numbers: []string{newNumber},
	//}

	//if err := svc.tclient.Send(ctx, Cases.tplId, Cases.params, Cases.numbers...); err != nil {
	//	//打个日志
	//	return err
	//}
	u.Verification = code
	err = svc.repo.SetVerification(ctx, u)
	return err

}
func (svc *UserService) Edit(ctx context.Context, phone string, op string, np string) error {
	//看看op 也就是原密码和数据库的密码是不是相同
	//查询原密码
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != nil {
		return err
	}
	//比较密码是否相同
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(op))
	if err != nil {
		//密码不对返回错误
		return ErrInvalidPassword
	}
	//修改新密码
	//生成新密码
	hash, err := bcrypt.GenerateFromPassword([]byte(np), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = svc.repo.Update(ctx, domain.User{
		Phone:    phone,
		Password: string(hash),
	})
	if err != nil {
		return err
	}
	return nil
}

// 登录
func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	//先查是否有这个用户
	u, err := svc.repo.FindByPhone(ctx, user.Phone)
	//没有这个用户
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		//返回账号或者密码错误
		return domain.User{}, err
	}
	return u, nil
}

// 注册
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//查看是否输入超过三次
	//比较验证码是否相同
	user, err1 := svc.repo.GetVerification(ctx, u)
	if err1 != nil {
		return err1
	}
	if user.Verification != u.Verification {
		return ErrVerificationFalse
	}
	//加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	if err := svc.repo.Create(ctx, u); err != nil {
		return err
	}
	if err := svc.repo.RemoveCode(ctx, u); err != nil {
		return err
	}
	return nil
}
func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
