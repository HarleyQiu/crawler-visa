package utils

import (
	"fmt"
	"time"
)

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
	return fmt.Sprintf("\n\n\n签证状态：%s\n创建日期：%s\n最后更新：%s\n详细信息：%s\n预约号：%s\n护照号：%s\n\n\n",
		status, createdAt.Format("2006年1月2日"), lastUpdatedAt.Format("2006年1月2日"), content, applicationID, passportNumber)
}

// FormatPassportStatus 构造一个显示护照状态的格式化消息。
// 它在标准化的消息格式中包含提供的状态内容和护照号码。
//
// 参数：
// -content（string）：要包含在消息中的状态描述。
// -passportNumber（string）：用于识别相关护照的护照号。
//
// 返回：
// -string：传达护照状态和号码的格式化消息。
func FormatPassportStatus(content, passportNumber string) string {
	return fmt.Sprintf("\n\n\n当前您的护照状态是：%s\n护照号：%s\n\n\n", content, passportNumber)
}
