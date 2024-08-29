package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

// QueryLoader 结构体定义，使用泛型 T
type QueryLoader[T any] struct {
	FilePath string
}

// NewQueryLoader 是 QueryLoader 的构造函数，返回一个 QueryLoader 的实例
func NewQueryLoader[T any](filePath string) *QueryLoader[T] {
	return &QueryLoader[T]{FilePath: filePath}
}

// LoadQueries 从JSON文件加载数据，返回切片类型为 T
func (ql *QueryLoader[T]) LoadQueries() ([]T, error) {
	file, err := os.Open(ql.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var queries []T
	if err := json.Unmarshal(bytes, &queries); err != nil {
		return nil, err
	}

	return queries, nil
}

// StatusTracker 类定义，使用泛型 T
type StatusTracker[T comparable] struct {
	statusMap map[string]T
	mu        sync.Mutex // 使用互斥锁保证并发安全
}

// NewStatusTracker 创建一个新的 StatusTracker 实例
func NewStatusTracker[T comparable]() *StatusTracker[T] {
	return &StatusTracker[T]{
		statusMap: make(map[string]T),
	}
}

// UpdateStatus 更新给定 ApplicationID 的状态，并返回是否有变化
func (st *StatusTracker[T]) UpdateStatus(key string, newStatus T) bool {
	st.mu.Lock()         // 修改map前加
	defer st.mu.Unlock() // 在方法结束时解锁

	currentStatus, exists := st.statusMap[key]
	if !exists || currentStatus != newStatus {
		st.statusMap[key] = newStatus // 更新状态
		return true                   // 状态改变或者是新的状态
	}
	return false // 没有变化
}

// NotificationSender 定义了用于发送通知的结构体，包含发送通知的 URL。
type NotificationSender struct {
	NotificationURL string // 发送通知的具体 URL。
}

// NotificationData 定义了通知所需的数据结构，包括系统类型、消费区域、监控国家等信息。
type NotificationData struct {
	Sys        string `json:"sys"`
	ConsDist   string `json:"consDist"`
	MonCountry string `json:"monCountry"`
	ApptTime   string `json:"apptTime"`
	Status     string `json:"status"`
	UserName   string `json:"userName"`
	Remark     string `json:"remark"`
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

// FormatVisaStatus 格式化签证状态信息并返回详细描述文本。
// 此函数接收签证状态（status）、详细信息（content）、创建日期（created）、
// 最后更新日期（lastUpdated）、预约号（applicationID）以及护照号（passportNumber）作为输入参数。
// 返回的字符串包含了所有这些信息，格式化后易于阅读。
//
// 参数:
//
//	status string - 签证的当前状态。
//	content string - 关于签证状态的附加信息。
//	created string - 签证创建的日期，格式应为 "02-Jan-2006"。
//	lastUpdated string - 签证最后更新的日期，格式应为 "02-Jan-2006"。
//	applicationID string - 签证的预约号。
//	passportNumber string - 护照号码。
//
// 返回值:
//
//	string - 格式化后的签证状态描述，包括所有提供的信息。
//
// 示例:
//
//	statusText := FormatVisaStatus("已批准", "请按时前往大使馆", "01-Jan-2023", "10-Jan-2023", "AB123456", "123456789")
//	fmt.Println(statusText)
//
// 输出将是:
//
//	签证状态：已批准
//	创建日期：2023年1月1日
//	最后更新：2023年1月10日
//	详细信息：请按时前往大使馆
//	预约号：AB123456
//	护照号：123456789
//
// 注意: 本函数不处理解析日期时的错误，调用者需确保提供的日期格式正确。
func FormatVisaStatus(status, content, created, lastUpdated, applicationID, passportNumber string) string {
	// 解析日期字符串
	createdAt, _ := time.Parse("02-Jan-2006", created)
	lastUpdatedAt, _ := time.Parse("02-Jan-2006", lastUpdated)

	// 组织成描述性文本，包括预约号和护照号
	return fmt.Sprintf("\n签证状态：%s\n创建日期：%s\n最后更新：%s\n详细信息：%s\n预约号：%s\n护照号：%s\n",
		status, createdAt.Format("2006年1月2日"), lastUpdatedAt.Format("2006年1月2日"), content, applicationID, passportNumber)
}
