# gotools

用 Python 的写法处理 Go 里的日常数据。

`gotools` 把 Python 3.13 常用的内置函数、类型方法和标准库模块，映射成**零第三方依赖**的 Go API：字符串、列表、字典、集合、JSON、CSV、文件/路径、时间、正则、日志、HTTP、子进程等。

目标不是嵌入 Python 解释器，而是在 Go 里保留 Python 的表达习惯，同时保留 Go 的性能和类型安全。

## 特点

- **零依赖**：只用 Go 标准库。
- **类型安全优先**：能用泛型的地方用泛型，无法静态表达时才提供动态版。
- **命名贴近 Python**：`sum` → `Sum`，`most_common` → `Most_common`。
- **诚实标注差异**：每个不能完全模拟的地方都有 `// DIFF:` 注释。

## 安装

```bash
go get github.com/BK-Go-Python/gotools
```

要求 **Go 1.26+**。

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/BK-Go-Python/gotools/list"
    "github.com/BK-Go-Python/gotools/py"
    "github.com/BK-Go-Python/gotools/set"
    "github.com/BK-Go-Python/gotools/str"
)

func main() {
    fmt.Println(str.Title("hello world"))                    // Hello World
    fmt.Println(str.Format("{} 有 {:.2f} 元", "Tom", 12.345)) // Tom 有 12.35 元
    fmt.Println(str.Split("a,b,c", ","))                     // [a b c]

    nums := list.Of(3, 1, 2)
    fmt.Println(py.Sorted(nums))                             // [1 2 3]
    fmt.Println(py.Sum(nums))                                // 6
    fmt.Println(py.Map(func(v int) int { return v * 2 }, nums)) // [6 2 4]

    s := set.New(1, 2, 3).Union(set.New(3, 4))
    fmt.Println(s.Len())                                     // 4

    smallest, err := py.Min([]int{3, 1, 2})
    fmt.Println(smallest, err)                               // 1 <nil>
}
```

> 所有包都以“首字母大写、其余保持 Python 拼写”为命名约定。`str`、`list`、`dict`、`set`、`string`、`json` 等包名与标准库或内建标识符接近，导入时可按需起别名。

## 按场景使用

### 字符串

```go
fmt.Println(str.Title("hello world"))          // Hello World
fmt.Println(str.Center("go", 6, '-'))          // --go--
fmt.Println(str.Removeprefix("prefix-x", "prefix-")) // x
fmt.Println(str.Split("a,b,c", ","))           // [a b c]
fmt.Println(str.Join("-", []string{"a", "b"})) // a-b
fmt.Println(str.Format("{} 有 {:.2f} 元", "Tom", 12.345))
```

字符串按 rune 处理，中文、emoji 与 Python 3 一致。

### 列表、字典、集合

```go
nums := list.New(3, 1, 2)
nums.Sort()
nums.Append(4).Insert(0, 0)
fmt.Println(nums.String()) // [0 1 2 3 4]

d := dict.New().Set("name", "Tom").Set("age", 18)
fmt.Println(d.Get("name")) // Tom
fmt.Println(d.Keys())      // [age name]

s := set.New(1, 2, 3).Union(set.New(3, 4))
fmt.Println(s.Len())       // 4
fmt.Println(s.Has(2))      // true
```

`list.List` 与 `dict.Dict` 是原地修改的 Python 风格容器；泛型 `set.Set[T]` 的 `Add` 返回新集合，常用 `set.New`。

### Python 内置函数

```go
fmt.Println(py.Range(1, 5)) // [1 2 3 4]

squares := py.Map(func(v int) int { return v * v }, []int{1, 2, 3})
fmt.Println(squares) // [1 4 9]

evens := py.Filter(func(v int) bool { return v%2 == 0 }, []int{1, 2, 3, 4})
fmt.Println(evens) // [2 4]

fmt.Println(py.Sum([]int{1, 2, 3})) // 6

pairs := py.Enumerate([]string{"a", "b"}, 1)
fmt.Println(pairs) // [{1 a} {2 b}]
```

泛型版保留元素类型：

```go
py.Sum([]int{1, 2, 3})       // int
py.Sum([]float64{1.5, 2.5})  // float64
py.Min([]int{3, 1, 2})       // (int, error)
```

需要处理 `[]any` 或 map 等异构输入时，使用 `SumAny`、`MinAny`、`MapAny` 等动态版。

### JSON 与 CSV

```go
value, err := json.Loads(`{"name":"Tom","age":18}`)
text, err := json.Dumps(value)
pretty, err := json.DumpsIndent(value, "  ")

_ = json.Dump("data.json", value)
value, err = json.Load("data.json")
```

```go
rows, err := csv.ReadFile("data.csv")
dicts, err := csv.ReadDictsFile("data.csv")

_ = csv.WriteFile("out.csv", [][]string{{"name", "age"}, {"Tom", "18"}})
_ = csv.WriteDictsFile("out.csv",
    []string{"name", "age"},
    []map[string]string{{"name": "Tom", "age": "18"}})
```

### 文件与路径

```go
p := pathlib.New("data", "users.json")
_ = p.WriteText(`{"ok":true}`)

text, err := p.ReadText()
exists := p.Exists()
items, err := p.Parent().Iterdir()
matches, err := p.Parent().Rglob("*.json")

fmt.Println(ospath.Join("data", "users.json"))
content, err := file.Read("data/users.json")
```

`pathlib` 是对象式 Path；`file` 是函数式文件操作；`ospath` 是 `os.path` 风格的字符串工具。

### 时间与日期

```go
import (
    "time"

    "github.com/BK-Go-Python/gotools/datetime"
    gotime "github.com/BK-Go-Python/gotools/time"
)

t, err := gotime.Strptime("2026-10-20 08:42:53", "%Y-%m-%d %H:%M:%S")
fmt.Println(gotime.Strftime(t, "%Y年%m月%d日 %H:%M"))

d, err := datetime.NewDate(2026, time.October, 20)
fmt.Println(d.AddDays(1)) // 2026-10-21
```

### 正则、文本与编码

```go
all, err := re.FindAll(`\d+`, "a1b22c333") // [1 22 333]
out, err := re.Sub(`(\w+)@(\w+)`, `${2}@${1}`, "tom@example")

fmt.Println(textwrap.Wrap("The quick brown fox", 10))
fmt.Println(base64.B64encode([]byte("hello"))) // aGVsbG8=
fmt.Println(html.Escape(`<b>"x"</b>`))
```

`re` 底层是 Go RE2：不支持 lookaround/回溯引用；替换模板用 `$1`，不是 `\1`。

### HTTP、环境与日志

```go
resp, err := httpclient.Get("https://example.com")
fmt.Println(resp.StatusCode, resp.Text())

var payload map[string]any
err = httpclient.GetJSON("https://example.com/data.json", &payload)

fmt.Println(env.Get("HOME", "/tmp"))
logger := logging.New("app")
logger.Info("started")
logger.Warningf("slow: %dms", 120)
```

```go
out, err := subprocess.CheckOutput("go", "version")
fmt.Println(out)

result, err := subprocess.Run("echo", "hello")
fmt.Println(result.Stdout)
```

### 集合类型、迭代器与缓存

```go
c := collections.NewCounter("a", "a", "b")
fmt.Println(c.Count("a"))       // 2
fmt.Println(c.Most_common(2))   // [{a 2} {b 1}]

d := collections.NewDeque(1, 2, 3)
d.AppendLeft(0)

fmt.Println(itertools.Chain([]int{1, 2}, []int{3, 4})) // [1 2 3 4]

square := functools.Memoize(func(n int) int { return n * n })
fmt.Println(square(5)) // 25

cache := functools.NewLruCache[string, int](2)
cache.Put("a", 1)
value, ok := cache.Get("a")
```

## 包索引

| 分类 | 包 | 说明 |
|---|---|---|
| 核心类型 | [`str`](str)、[`list`](list)、[`dict`](dict)、[`set`](set)、[`bytearray`](bytearray)、[`conv`](conv) | 字符串、列表、字典、集合、可变字节、类型转换 |
| 内置函数 | [`py`](py) | `range/zip/enumerate/map/filter/sorted/sum/min/max/print/input` 等 |
| 容器与迭代 | [`collections`](collections)、[`itertools`](itertools)、[`functools`](functools) | Counter/Deque/OrderedDict/ChainMap、组合迭代器、缓存/Partial |
| 数据格式 | [`json`](json)、[`csv`](csv) | loads/dumps、CSV 读写与 DictReader/DictWriter |
| 文件与路径 | [`file`](file)、[`pathlib`](pathlib)、[`ospath`](ospath)、[`shutil`](shutil) | 函数式文件操作、Path 对象、os.path 字符串工具、文件树复制 |
| 时间 | [`time`](time)、[`datetime`](datetime) | 时间戳/格式化/日期计算、Date/DateTime/Timedelta |
| 文本与编码 | [`re`](re)、[`textwrap`](textwrap)、[`string`](string)、[`html`](html)、[`base64`](base64)、[`urllib`](urllib) | 正则、换行、模板、HTML 转义、Base64、URL |
| 网络与进程 | [`httpclient`](httpclient)、[`subprocess`](subprocess)、[`signal`](signal) | HTTP 请求、命令执行、信号 |
| 系统与运行时 | [`sys`](sys)、[`env`](env)、[`logging`](logging) | argv/进程信息、环境变量、日志 |
| 数值与随机 | [`statistics`](statistics)、[`random`](random)、[`hash`](hash) | 统计、随机数、哈希/UUID |

## 命名约定

Go 的导出标识符必须以大写字母开头，所以无法提供真正小写的 `sum`/`most_common`。约定是：

- **只把首字母大写，其余保持 Python 拼写**：`takewhile` → `Takewhile`，`most_common` → `Most_common`，`bit_count` → `Bit_count`。
- **泛型优先**：`py.Sum([]int{...})` 直接返回 `int`；动态版是 `SumAny`。
- **旧名字保留为 Deprecated 别名**：如 `SumOf`、`MinOf`、`MapOf`，避免已有代码断掉。
- **需要多个返回值的语义用 `error`**：越界、IO、解析失败不 panic。

## 兼容性说明

这是一个“Python 风格”库，不是 Python 运行时：

- 不实现 `eval/exec`、描述符、类继承体系、异步协议等解释器能力。
- `re` 使用 Go RE2，不支持 lookaround/回溯引用。
- `datetime` 没有 Python `timezone/tzinfo` 对象，用 `time.Location` 表达。
- `json` 的大整数保留为 `json.Number`，不是任意精度整数。
- `subprocess` 没有 `Popen` 流式管道；复杂进程控制请直接用 `os/exec`。
- `logging` 没有 handler/formatter/logger 层级注册表。

更详细的逐项差异：

- 类型方法：[`docs/type-methods.md`](docs/type-methods.md)
- 不可模拟的内置函数：[`docs/unsupported-interpreter.md`](docs/unsupported-interpreter.md)、[`docs/unsupported-protocol.md`](docs/unsupported-protocol.md)
- 各模块文档：[`docs/`](docs)

## 测试

```bash
go test ./...
go test -race ./...
go vet ./...
gofmt -l $(git ls-files '*.go')
```

测试期望值尽量以 Python 3.13 实际行为为准；新增 API 使用表驱动测试，覆盖正常、边界和错误路径。

## License

CC0 1.0 Universal，见 [LICENSE](LICENSE)。
