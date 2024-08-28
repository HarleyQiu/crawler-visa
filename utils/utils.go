package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
