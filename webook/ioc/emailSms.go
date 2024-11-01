package ioc

import (
	"fmt"
	"go-study/webook/internal/service/sms/email"
	"gopkg.in/gomail.v2"
	"os"
)

func InitEmailSms() *email.Service {
	//初始化邮箱验证码配置
	m := gomail.NewMessage()
	m.SetHeader("From", "HuaHua<"+"2041436630@qq.com"+">") // 设置发件人别名
	m.SetHeader("Subject", "您的验证码")                        // 邮件主题

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
	el := email.NewService(d, m)
	return el
}
