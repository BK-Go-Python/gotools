// Package str 提供 Python str 风格的方法：大小写转换、裁剪、分割、查找、
// 判断、对齐填充、切片以及 str.format 风格的字符串格式化。
//
// 所有函数都以 rune（字符）为单位处理字符串，因此中文、emoji 等
// 多字节字符的行为与 Python 3 一致，而不是按字节处理。
//
// 该包不依赖任何第三方库。
package str

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Str 把任意值转换为字符串。
//
// 数值按十进制输出，字符串原样返回，[]byte 转成 string；其他类型会尝试
// 序列化为 JSON，失败时退化为 fmt.Sprint。nil 返回空字符串。
//
// 这是项目最早提供的函数，行为偏 Go 习惯；如果需要 Python str() 的
// 语义（None、True/False、带引号的容器），请使用 py.Str。
func Str(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		return string(data)
	}
}
