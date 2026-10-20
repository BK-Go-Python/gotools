# gotools

用 Python 的方式写 Go。

`gotools` 把 Python 3.13 常用的内置函数与类型方法，按 Python 同名 API 提供 Go
实现：字符串、列表、字典、集合、类型转换、时间、哈希、文件操作，以及
`range/zip/enumerate/map/filter/sorted/sum/min/max` 等内置函数。

目标不是嵌入 Python 解释器，而是让 Go 在常见数据处理上拥有接近 Python 的
表达力，同时保持 Go 的性能与类型安全。

## 安装

```bash
go get github.com/BK-Go-Python/gotools
```

要求 **Go 1.26+**（使用泛型、`cmp`、`slices`、`maps`、`iter`），零第三方依赖。

## 包一览

| 包 | 内容 |
|---|---|
| `str` | Python `str` 方法：大小写、strip、split、find、startswith、center、zfill、切片、`str.format` 模板 |
| `list` | `list.List` 及 append/extend/insert/remove/pop/index/sort/reverse/slice 等方法 |
| `dict` | `dict.Dict` 及 get/setdefault/keys/values/items/update/pop/getPath 等方法 |
| `set` | `set.Set`、`set.Frozen`（`set` / `frozenset`） |
| `collections` | `Counter`、`Deque`、`OrderedDict`（Python `collections` 风格） |
| `itertools` | `Chain`、`Islice`、`Repeat`、`Cycle`、`Product`、`Combinations`、`Permutations`、`Groupby` 等 |
| `py` | Python 内置函数：len/range/zip/enumerate/map/filter/sorted/sum/min/max/abs/round/format/hash/print/input 等 |
| `conv` | `int/float/bool/str/bytes/bytearray/complex` 转换、JSON、struct ↔ map |
| `time` | 时间戳、`strftime/strptime`、日期加减、月初月末、中文相对时间 |
| `hash` | MD5/SHA/HMAC、文件哈希、随机 token、UUID |
| `file` | os/os.path/shutil 风格文件操作与 `file.Open` |

`internal/pyval` 是内部值语义层，不对外使用。

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/BK-Go-Python/gotools/py"
    "github.com/BK-Go-Python/gotools/set"
    "github.com/BK-Go-Python/gotools/str"
)

func main() {
    fmt.Println(str.Title("hello world"))          // Hello World
    fmt.Println(str.Format("{} 有 {:.2f} 元", "Tom", 12.345)) // Tom 有 12.35 元
    fmt.Println(str.Split("a,b,c", ","))           // [a b c]

    fmt.Println(py.Range(5))                       // [0 1 2 3 4]
    fmt.Println(py.MapOf([]int{1, 2, 3}, func(v int) int { return v * 2 })) // [2 4 6]
    fmt.Println(py.SumOf([]int{1, 2, 3}))          // 6

    s := set.New(1, 2, 3).Union(set.New(3, 4))
    fmt.Println(s.Len())                           // 4

    fmt.Println(py.Format(3.14159, ".2f"))         // 3.14
}
```

类型转换与结构体映射：

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}

m, _ := conv.StructToMap(User{Name: "Tom", Age: 18})
// map[age:18 name:Tom]
var u User
_ = conv.MapToStruct(m, &u)
```

时间：

```go
import gotime "github.com/BK-Go-Python/gotools/time"

t, _ := gotime.Strptime("2026-10-20 08:42:53", "%Y-%m-%d %H:%M:%S")
fmt.Println(gotime.Strftime(t, "%Y年%m月%d日 %H:%M"))
```

## Python 3.13 内置函数覆盖

已按 Python 3.13 文档的 71 个内置函数做过映射与分片：

- **12 个可完全模拟**：如 `all`、`any`、`ascii`、`bin`、`bool`、`float`、`hex`、
  `len`、`oct`、`ord`、`print`。
- **42 个部分模拟**：如 `list`、`dict`、`set`、`sorted`、`sum`、`zip`、`open`，
  差异写在 GoDoc 的 `// DIFF:` 注释中。
- **17 个不提供**：依赖 CPython 解释器、异步协议或类/描述符模型，如 `eval`、
  `exec`、`globals`、`aiter`、`property`、`super`。

不可模拟项不导出空壳 API，避免永久占位符号；说明见：

- [`docs/unsupported-interpreter.md`](docs/unsupported-interpreter.md)
- [`docs/unsupported-protocol.md`](docs/unsupported-protocol.md)

常用类型方法（str/list/dict/set/int/float/bytes）的逐项对照见
[`docs/type-methods.md`](docs/type-methods.md)。

### 命名约定

Go 的导出标识符必须以大写字母开头，因此无法提供真正小写的 `sum`/`most_common`。
gotools 的约定是：**只把首字母大写，其余保持 Python 拼写**（`chain`→`Chain`、
`takewhile`→`Takewhile`、`most_common`→`Most_common`、`bit_count`→`Bit_count`、
`bytes.upper`→`Bytes_upper`）。已经存在的驼峰 API 保留为兼容别名，新 API 优先使用
Python 拼写。

- Python 名已被动态版占用时，泛型版本加 `Of` 后缀：
  `SumOf/MinOf/MaxOf/SortedOf/MapOf/FilterOf/EnumerateOf/ReversedOf/DivModOf`。
- 变参类型参数无法表达时用固定元数：`Zip2`/`Zip3`。
- 动态版保留 Python 语义；泛型版是类型安全的快速路径，不改变旧 API。

## 设计原则

1. **只增不改**：现有导出 API 的签名与语义保持兼容。
2. **新 API 不 panic**：解析、IO、反射、越界错误统一用 `error` 返回；
   历史动态 API 的 panic 行为仅为兼容保留。
3. **确定性**：map 迭代按键排序，排序稳定，避免 Go map 随机顺序影响结果。
4. **零依赖**：只用标准库；Go 1.26 的 `slices`、`maps`、`cmp`、`iter` 优先。
5. **诚实标注差异**：无法完全模拟 Python 的地方用 `// DIFF:` 写明，不假装兼容。
6. **不做解释器**：不实现 `eval/exec`，不通过 `os/exec` 调 Python 冒充兼容层。

## 测试

```bash
go test ./...
go test -race ./...
go vet ./...
gofmt -l str list dict set py conv time hash file internal
```

测试期望值以 Python 3.13 实际行为为准；新增函数使用表驱动测试，覆盖正常、边界和
错误路径。

## 仓库结构

```
gotools/
├── str/ list/ dict/ set/ collections/ itertools/ py/ conv/ time/ hash/ file/
├── internal/pyval/      # 内部值语义（truthiness/len/iter/compare/repr）
├── docs/                # 不可模拟函数的说明
├── go.mod
└── README.md
```


## License

CC0 1.0 Universal，见 [LICENSE](LICENSE)。
