// Package env 提供 Python os.environ 风格的环境变量工具。
//
// 所有函数都是对 os.LookupEnv/Setenv/Unsetenv 的薄封装；All 按 KEY=VALUE
// 排序返回，保证确定性。
package env

import (
	"os"
	"sort"
	"strings"
)

// Get 等价于 Python 的 os.environ.get(key, default)。
// 未提供 default 且变量不存在时返回空字符串。
func Get(key string, def ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// Lookup 等价于 Python 的 key in os.environ + 取值。
func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Set 等价于 Python 的 os.environ[key] = value。
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset 等价于 Python 的 del os.environ[key]。
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Has 等价于 Python 的 key in os.environ。
func Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Keys 返回所有环境变量名，按字典序排序。
func Keys() []string {
	entries := os.Environ()
	keys := make([]string, 0, len(entries))
	for _, entry := range entries {
		if idx := strings.IndexByte(entry, '='); idx >= 0 {
			keys = append(keys, entry[:idx])
		}
	}
	sort.Strings(keys)
	return keys
}

// All 返回所有环境变量，格式为 KEY=VALUE，按 KEY 排序。
func All() []string {
	entries := os.Environ()
	sort.Strings(entries)
	return entries
}

// Expand 等价于 Python 的 os.path.expandvars / os.environ 展开：
// 支持 $KEY 与 ${KEY}，未定义的变量保留原样（与 os.ExpandEnv 一致）。
func Expand(s string) string {
	return os.ExpandEnv(s)
}
