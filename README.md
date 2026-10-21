# gotools

用 Python 的写法处理 Go 里的日常数据。

`gotools` 把 Python 3.13 常用的内置函数、类型方法和标准库模块，映射成**零第三方依赖**的 Go API：字符串、列表、字典、集合、JSON、CSV、文件/路径、时间、正则、日志、HTTP、子进程等。

目标不是嵌入 Python 解释器，而是在 Go 里保留 Python 的表达习惯，同时保留 Go 的性能和类型安全。

## 安装

```bash
go get github.com/BK-Go-Python/gotools
```

要求 **Go 1.26+**。

## Go Notebook 快速开始

下面每个代码块都是 notebook 单元：**直接调用，不写 `func main()`**；`// =>` 后是返回值/输出。实际使用时补上对应包导入即可。

```go
str.Title("hello world")                            // => "Hello World"
str.Format("{} 有 {:.2f} 元", "Tom", 12.345)        // => "Tom 有 12.35 元"
str.Split("a,b,c", ",")                             // => [a b c]

nums := []int{3, 1, 2}
py.Sorted(nums)                                     // => [1 2 3]
py.Sum(nums)                                        // => 6
py.Map(func(v int) int { return v * 2 }, nums)      // => [6 2 4]

set.New(1, 2, 3).Union(set.New(3, 4)).Len()         // => 4
```

各包主要 API 的用法见：

- [`docs/notebook.md`](docs/notebook.md) —— 各包主要 API 的 notebook 单元（完整导出列表用 `go doc ./包名`）
- [`docs/type-methods.md`](docs/type-methods.md) —— str/list/dict/set/int/float/bytes 逐项对照

## 常用示例

### 字符串与容器

```go
str.Center("go", 6, '-')                            // => "--go--"
str.Removeprefix("prefix-x", "prefix-")             // => "x"
str.Join("-", []string{"a", "b"})                   // => "a-b"

nums := list.New(3, 1, 2)
nums.Sort()                                         // => [1, 2, 3]
nums.String()                                       // => "[1, 2, 3]"

d := dict.New().Set("name", "Tom").Set("age", 18)
d.Get("name")                                       // => "Tom"
d.Keys()                                            // => [age name]

s := set.New(1, 2, 3)
s.Union(set.New(3, 4)).Len()                        // => 4
s.Has(2)                                            // => true
```

### JSON、CSV、路径

```go
value, err := json.Loads(`{"name":"Tom"}`)
value.(map[string]any)["name"]                      // => "Tom"
err                                                 // => nil

json.Dumps(map[string]any{"name": "Tom"})           // => "{\"name\":\"Tom\"}", nil

csv.Read(strings.NewReader("a,b\n1,2\n"))          // => [[a b] [1 2]], nil

p := pathlib.New("data", "users.json")
p.Name()                                            // => "users.json"
p.Suffix()                                          // => ".json"
p.Parent().String()                                 // => "data"

ospath.Join("data", "users.json")                   // => "data/users.json"
```

### 时间、正则、统计

```go
t, _ := gotime.Strptime("2026-10-20 08:42:53", "%Y-%m-%d %H:%M:%S")
gotime.Strftime(t, "%Y年%m月%d日 %H:%M")             // => "2026年10月20日 08:42"

d, _ := datetime.NewDate(2026, time.October, 20)
d.AddDays(1).ISOFormat()                            // => "2026-10-21"

re.FindAll(`\d+`, "a1b22c333")                      // => [1 22 333], nil
re.Sub(`(\w+)@(\w+)`, `${2}@${1}`, "tom@example")   // => "example@tom", nil

statistics.Mean([]int{1, 2, 3})                     // => 2, nil
statistics.Quantiles([]int{1, 2, 3, 4}, 4)          // => [1.25 2.5 3.75], nil
```

### HTTP、进程、环境、日志

```go
resp, err := httpclient.Get("https://example.com")
resp.StatusCode                                     // => 200
resp.Text()                                         // => HTML 文本
err                                                 // => nil

out, err := subprocess.CheckOutput("go", "version")
out                                                 // => "go version go1.26.x ...\n"
err                                                 // => nil

env.Get("HOME", "/tmp")                             // => 当前 HOME
logger := logging.New("app")
logger.Info("started")                              // 输出: 2026-... [INFO] app started
```

### 集合类型、迭代器、缓存

```go
c := collections.NewCounter("a", "a", "b")
c.Count("a")                                        // => 2
c.Most_common(2)                                    // => [{a 2} {b 1}]

dq := collections.NewDeque(1, 2, 3)
dq.AppendLeft(0)
dq.Items()                                          // => [0 1 2 3]

itertools.Chain([]int{1, 2}, []int{3, 4})           // => [1 2 3 4]
itertools.Combinations([]int{1, 2, 3}, 2)           // => [[1 2] [1 3] [2 3]]

square := functools.Memoize(func(n int) int { return n * n })
square(5)                                           // => 25
```

> Notebook 片段不是完整程序；真实代码中需要写 `import`，并把表达式放进函数或 `main`。完整可运行示例见各包测试。

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
- `re` 使用 Go RE2，不支持 lookaround/回溯引用；替换模板用 `$1`。
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
