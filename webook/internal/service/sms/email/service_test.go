package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"testing"
)

func TestSend(t *testing.T) {
	//初始化邮箱验证码配置
	m := gomail.NewMessage()
	m.SetHeader("From", "HuaHua<"+"2041436630@qq.com"+">") // 设置发件人别名
	m.SetHeader("Subject", "您的验证码")
	// 邮件主题
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
	gomail1 := NewService(d, m)
	//要发送给谁
	testCases := []struct {
		name  string
		email string
	}{{name: "wang", email: "3504974916@qq.com"},
		{name: "wang", email: "3504974916@qq.com"},
		{name: "wang", email: "3504974916@qq.com"},
		{name: "wang", email: "3504974916@qq.com"},
		{name: "wang", email: "3504974916@qq.com"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gomail1.Send(tc.email)
		})
	}

}
