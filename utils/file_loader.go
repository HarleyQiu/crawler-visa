package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// QueryLoader 定义了一个泛型结构体，用于从JSON文件中加载数据。
// 使用泛型 T 允许 QueryLoader 加载任何类型的数据。
type QueryLoader[T any] struct {
	FilePath string
}

// NewQueryLoader 创建并返回一个新的 QueryLoader 实例，初始化文件路径。
// 参数:
//
//	filePath string - 要加载的JSON文件的路径。
//
// 返回值:
//
//	*QueryLoader[T] - 新创建的 QueryLoader 实例。
func NewQueryLoader[T any](filePath string) *QueryLoader[T] {
	return &QueryLoader[T]{FilePath: filePath}
}

// LoadQueries 从设置的文件路径中加载JSON文件，并将其解析为切片类型为 T 的数据。
// 返回值包括加载的数据切片和可能发生的错误。
// 如果文件打开或读取失败，或JSON解析不成功，将返回错误。
//
// 返回值:
//
//	[]T - 加载并解析的数据。
//	error - 在加载或解析过程中遇到的错误。
//
// 示例:
//
//	loader := NewQueryLoader[MyDataType]("path/to/data.json")
//	data, err := loader.LoadQueries()
//	if err != nil {
//	    fmt.Println("加载查询失败:", err)
//	} else {
//	    fmt.Println("加载的数据:", data)
//	}
//
// 注意: 此函数未处理特定的JSON格式问题，使用前应确保JSON文件格式与泛型类型T匹配。
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
