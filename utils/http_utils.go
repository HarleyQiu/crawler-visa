package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ParseBody 从 HTTP 请求中读取并解析 JSON 格式的请求体。
// 参数:
//
//	r *http.Request - HTTP 请求指针，包含待解析的请求体。
//	x interface{} - 一个指向将要存储解析数据的变量的指针。
//	该变量应预定义为适当的结构体或其他类型，以便能够正确解析 JSON 数据。
//
// 无返回值。如果解析过程中发生错误，例如读取请求体或解析 JSON 时，该函数将不会对输入的变量 x 进行修改。
//
// 示例:
//
//	type UserData struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	var data UserData
//	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name":"John", "age":30}`))
//	ParseBody(r, &data)
//	fmt.Printf("Parsed data: %+v\n", data)
//
// 注意: 该函数不返回错误，调用者无法知道解析是否成功。
// 调用者应当通过检查解析后变量的状态或设计函数以返回错误信息，以便更好地处理错误。
func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
