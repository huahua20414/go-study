package email

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"strconv"
	"time"
)

type Service struct {
	d *gomail.Dialer
	m *gomail.Message
}

func NewService(d *gomail.Dialer, m *gomail.Message) *Service {
	return &Service{d: d, m: m}
}
func (svc *Service) Send(toname string) string {
	rand.Seed(time.Now().UnixNano())     // 以当前时间为随机数种子
	number := rand.Intn(900000) + 100000 // 生成100000到999999之间的随机数
	code := strconv.Itoa(number)         // 生成六位数验证码

	message := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <title>邮件验证码模板</title>
      <base target="_blank"/>
      <style type="text/css">
        ::-webkit-scrollbar {
          display: none;
        }
      </style>
    </head>
    <body>
    <table width="700" border="0" align="center" cellspacing="0">
      <tbody>
      <tr>
        <td>
          <div style="width:680px;padding:0 10px;margin:0 auto;">
            <div style="line-height:1.5;font-size:14px;margin-bottom:25px;color:#4d4d4d;">
              <strong style="display:block;margin-bottom:15px;">尊敬的用户：<span
                style="color:#f60;font-size: 16px;"></span>您好！</strong>
              <strong style="display:block;margin-bottom:15px;">
                这是您的<span style="color: red">验证码</span>，请在验证码输入框中输入：<span
                style="color:#f60;font-size: 24px">%s</span>，以完成操作。
              </strong>
            </div>
            <div style="margin-bottom:30px;">
              <small style="display:block;margin-bottom:20px;font-size:12px;">
                <p style="color:#747474;">
                  注意：此操作可能会修改您的密码、登录邮箱或绑定手机。如非本人操作，请及时登录并修改密码以保证帐户安全
                  <br>（工作人员不会向你索取此验证码，请勿泄漏！)
                </p>
              </small>
            </div>
          </div>
        </td>
      </tr>
      </tbody>
    </table>
    </body>
    </html>
    `, code)

	svc.m.SetHeader("To", toname)                               // 收件人
	svc.m.SetHeader("Content-Type", "text/html; charset=UTF-8") // 设置内容类型和字符编码
	svc.m.SetBody("text/html", message)                         // 邮件内容
	// 关闭SSL协议认证
	svc.d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := svc.d.DialAndSend(svc.m); err != nil {
		panic(err)
	}
	return code
}
