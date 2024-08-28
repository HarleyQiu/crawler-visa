package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
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
