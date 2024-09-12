package utils

import (
	"encoding/json"
	"net/http"
)

// ResultData 定义统一的响应结构
type ResultData struct {
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
	Code    int         `json:"code"`
}

func ResultJSON(w http.ResponseWriter, data interface{}, message string, optionalStatus ...int) {
	status := http.StatusOK // 默认状态码为 200 OK
	if len(optionalStatus) > 0 {
		status = optionalStatus[0] // 如果提供了状态码，则使用提供的状态码
	}

	// 检查 data 是否为 string 类型
	if jsonData, ok := data.(string); ok {
		var parsedData interface{}
		err := json.Unmarshal([]byte(jsonData), &parsedData)
		if err == nil {
			data = parsedData // 如果是有效的 JSON 字符串，使用解析后的对象
		}
		// 如果不是有效的 JSON 字符串，data 将保持原样（可能会返回原始字符串）
	}

	response := ResultData{
		Code:    status,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ResultError 用于发送错误响应
func ResultError(w http.ResponseWriter, message string, status int) {
	ResultJSON(w, nil, message, status)
}
