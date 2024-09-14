package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NotificationSender 定义了用于发送通知的结构体，包含发送通知的 URL。
type NotificationSender struct {
	NotificationURL string // 发送通知的具体 URL。
}

// NotificationData 定义了通知所需的数据结构，包括系统类型、消费区域、监控国家等信息。
type NotificationData struct {
	Sys        string `json:"sys"`        // 系统标识
	ConsDist   string `json:"consDist"`   // 领区
	MonCountry string `json:"monCountry"` // 国家
	ApptTime   string `json:"apptTime"`   // 预约时间
	Status     string `json:"status"`     // 状态
	UserName   string `json:"userName"`   // 用户名
	Remark     string `json:"remark"`     // 备注
}

// NewNotificationSender 初始化一个新的 NotificationSender 实例。
// 参数:
//
//	notificationURL string - 用于发送通知的 URL。
//
// 返回值:
//
//	*NotificationSender - 新创建的 NotificationSender 实例。
func NewNotificationSender(notificationURL string) *NotificationSender {
	return &NotificationSender{NotificationURL: notificationURL}
}

// SendNotification 发送一个包含指定数据的 HTTP POST 请求到预设的 URL。
// 参数:
//
//	data NotificationData - 要发送的通知数据。
//
// 返回值:
//
//	error - 发送失败时返回错误。
//
// 示例:
//
//	sender := NewNotificationSender("https://example.com/notify")
//	data := NotificationData{
//	  Sys: "system1", ConsDist: "Area51", MonCountry: "USA",
//	  ApptTime: "12:00", Status: "Pending", UserName: "user123",
//	  Remark: "Urgent"
//	}
//	err := sender.SendNotification(data)
//	if err != nil {
//	    fmt.Println("通知发送失败:", err)
//	}
//
// 注意: 本函数处理了请求和响应的基本逻辑，但未实现错误重试机制或响应结果的复杂处理。
func (ns *NotificationSender) SendNotification(data NotificationData) error {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", ns.NotificationURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 打印响应体内容到控制台
	fmt.Printf("来自服务器的响应: %s\n", string(responseBody))
	return nil
}
