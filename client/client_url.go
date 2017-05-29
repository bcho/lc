package client

import "fmt"

func UrlSignin() string {
	return "https://leancloud.cn/1/signin"
}

func UrlSMSRecords(appId string) string {
	return fmt.Sprintf("https://leancloud.cn/1/clients/self/apps/%s/sms", appId)
}
