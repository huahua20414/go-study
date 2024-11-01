package ioc

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go-study/webook/internal/service/sms/tencent"
	"os"
)

func InitTencentSms() *tencent.Service {
	//初始化腾讯云短信配置
	//创建sms对象
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		fmt.Println("SMS_SECRET_ID is empty")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		fmt.Println("SMS_SECRET_KEY is empty")
	}
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		//打印到日志中
	}
	//模板id，签名名称,创建sms接口对象
	s := tencent.NewService(c, "1400946141", "goLang科技公众号")
	return s
}
