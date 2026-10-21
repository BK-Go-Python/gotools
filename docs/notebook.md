# Go Notebook：gotools 用法

每个代码块都是一个 notebook 单元：**直接调用，不写 `func main()`**；`// =>` 后面是返回值或输出。实际使用时补上对应包导入即可。

> 覆盖各包的主要 API；同一个函数的 Deprecated 别名、`Of` 兼容名不重复列出。完整签名用 `go doc ./包名` 查看。

## py：内置函数

```go
py.Sum([]int{1, 2, 3})                              // => 6
py.Min([]int{3, 1, 2})                              // => 1, nil
py.Max([]int{3, 1, 2})                              // => 3, nil
py.Sorted([]int{3, 1, 2})                           // => [1 2 3]

py.Map(func(v int) int { return v * 2 }, []int{1, 2, 3}) // => [2 4 6]
py.Filter(func(v int) bool { return v%2 == 0 }, []int{1, 2, 3, 4}) // => [2 4]
py.Enumerate([]string{"a", "b"}, 1)                 // => [{1 a} {2 b}]
py.Reversed([]int{1, 2, 3})                         // => [3 2 1]
py.Range(1, 5)                                      // => [1 2 3 4]
py.Zip2([]int{1, 2}, []string{"a", "b"})            // => [{1 a} {2 b}]

py.Len([]int{1, 2})                                 // => 2
py.Abs(-3)                                          // => 3
py.Round(2.675, 2)                                  // => 2.67
py.Bin(5)                                           // => "0b101"
py.Hex(255)                                         // => "0xff"
py.Oct(8)                                           // => "0o10"
str.Chr(65)                                         // => "A", nil
str.Ord("A")                                        // => 65, nil
py.Divmod(-7, 3)                                    // => -3, 2, nil
py.BitLength(8)                                     // => 4
py.BitCount(7)                                      // => 3

py.All([]int{1, 2})                                 // => true
py.Any([]int{0, 1})                                 // => true
py.Bool(0)                                          // => false
py.Format(3.14159, ".2f")                           // => "3.14", nil
```

## str：字符串

```go
str.Title("hello world")                            // => "Hello World"
str.Upper("abc")                                    // => "ABC"
str.Lower("ABC")                                    // => "abc"
str.Casefold("ß")                                   // => "ss"
str.Strip("  x  ")                                  // => "x"
str.Removeprefix("prefix-x", "prefix-")             // => "x"
str.Removesuffix("x-suffix", "-suffix")             // => "x"

str.Split("a,b,c", ",")                             // => [a b c]
str.Join("-", []string{"a", "b"})                   // => "a-b"
str.Replace("aaa", "a", "b", 2)                     // => "bba"
str.Center("go", 6, '-')                            // => "--go--"
str.Format("{} 有 {:.2f} 元", "Tom", 12.345)        // => "Tom 有 12.35 元"

str.Find("banana", "na")                            // => 2
str.Count("banana", "na")                           // => 2
str.Startswith("hello", "he")                       // => true
str.Endswith("hello", "lo")                         // => true
str.Isdigit("123")                                  // => true
str.Isalpha("abc")                                  // => true
```

## list / dict / set

```go
nums := list.New(3, 1, 2)
nums.Sort()                                         // => [1, 2, 3]
nums.String()                                       // => "[1, 2, 3]"
nums.Contains(2)                                    // => true
nums.Slice(0, 2).String()                           // => "[1, 2]"

list.Of(1, 2, 3)                                    // => []int{1, 2, 3}
```

```go
d := dict.New().Set("name", "Tom").Set("age", 18)
d.Get("name")                                       // => "Tom"
d.Keys()                                            // => [age name]
d.Items()                                           // => [{age 18} {name Tom}]
d.Has("name")                                       // => true
d.Merge(dict.New().Set("city", "SH")).Len()         // => 3
```

```go
a := set.New(1, 2, 3)
a.Has(2)                                            // => true
a.Union(set.New(3, 4)).Len()                        // => 4
a.Intersection(set.New(2, 3, 4)).Len()              // => 2
a.Difference(set.New(3)).Len()                      // => 2
set.New(1, 2).IsSubset(a)                           // => true
set.New(1, 2, 3, 4).Issuperset(a)                  // => true
```

## collections

```go
c := collections.NewCounter("a", "a", "b")
c.Count("a")                                        // => 2
c.Total()                                           // => 3
c.Most_common(2)                                    // => [{a 2} {b 1}]

dq := collections.NewDeque(1, 2, 3)
dq.AppendLeft(0)
dq.Items()                                          // => [0 1 2 3]
dq.PopLeft()                                        // => 0, true

od := collections.NewOrderedDict[string, int]()
od.Set("a", 1).Set("b", 2)
od.Keys()                                           // => [a b]
od.Popitem()                                        // => {b 2}, true

cm := collections.NewChainMap(map[string]int{"a": 1}, map[string]int{"b": 2})
cm.Get("b")                                         // => 2

dd := collections.NewDefaultDict[string](func() int { return 7 })
dd.Get("x")                                         // => 7
```

## bytearray

```go
b := bytearray.New('a', 'b')
b.Append('c').Insert(0, 'x')
string(b)                                           // => "xabc"
string(b.Upper())                                   // => "XABC"
b.Count('a')                                        // => 1
b.Replace([]byte("a"), []byte("A"))
string(b)                                           // => "xAbc"
```

## itertools

```go
itertools.Chain([]int{1, 2}, []int{3, 4})           // => [1 2 3 4]
itertools.Repeat("x", 3)                            // => [x x x]
itertools.Islice([]int{0, 1, 2, 3}, 1, 4, 2)        // => [1 3], nil
itertools.Product([]int{1, 2}, []int{3, 4})         // => [[1 3] [1 4] [2 3] [2 4]]
itertools.Combinations([]int{1, 2, 3}, 2)           // => [[1 2] [1 3] [2 3]]
itertools.Permutations([]int{1, 2}, 2)              // => [[1 2] [2 1]]
itertools.Combinations_with_replacement([]int{1, 2}, 2) // => [[1 1] [1 2] [2 2]]
itertools.Filterfalse([]int{1, 2, 3}, func(v int) bool { return v%2 == 1 }) // => [2]
itertools.Compress([]string{"a", "b", "c"}, []bool{true, false, true}) // => [a c]
itertools.Tee([]int{1, 2}, 2)                       // => [[1 2] [1 2]]
```

## functools

```go
square := functools.Memoize(func(n int) int { return n * n })
square(4)                                           // => 16
square(4)                                           // => 16（第二次命中缓存）

cache := functools.NewLruCache[string, int](2)
cache.Put("a", 1)
cache.Get("a")                                      // => 1, true
cache.Keys()                                        // => [a]

inc := functools.Partial2(func(a, b int) int { return a + b }, 1)
inc(2)                                              // => 3

partial, _ := functools.NewPartial(fmt.Sprintf, "n=%d")
partial.Call(3)                                     // => ["n=3"], nil
```

## json / csv

```go
json.Loads(`{"name":"Tom","age":18}`)               // => map[string]interface{}{...}, nil；取 ["name"] 得到 "Tom"
json.Dumps(map[string]any{"name": "Tom"})           // => "{\"name\":\"Tom\"}", nil
json.DumpsIndent(map[string]any{"a": 1}, "  ")      // => "{\n  \"a\": 1\n}", nil
json.Indent(`{"a":1}`, "  ")                        // => "{\n  \"a\": 1\n}", nil
```

```go
csv.Read(strings.NewReader("a,b\n1,2\n"))          // => [[a b] [1 2]], nil
csv.ReadDicts(strings.NewReader("a,b\n1,2\n"))     // => [map[a:1 b:2]], nil
var buf bytes.Buffer
csv.Write(&buf, [][]string{{"a", "b"}})             // => nil；buf.String() == "a,b\r\n"
csv.WriteDicts(&buf, []string{"a"}, []map[string]string{{"a": "1"}}) // => nil
```

## pathlib / ospath / file

```go
p := pathlib.New("data", "users.json")
p.Name()                                            // => "users.json"
p.Suffix()                                          // => ".json"
p.Parent().String()                                 // => "data"
p.Exists()                                          // => false（路径不存在时）

p.WriteText(`{"ok":true}`)                          // => nil
p.ReadText()                                        // => "{\"ok\":true}", nil
p.Resolve()                                         // => Path("/.../data/users.json"), nil
p.Parent().Rglob("*.json")                          // => []Path, nil
```

```go
ospath.Join("data", "users.json")                   // => "data/users.json"
ospath.Base("/a/b.txt")                             // => "b.txt"
ospath.Dir("/a/b.txt")                              // => "/a"
ospath.Ext("/a/b.txt")                              // => ".txt"
ospath.Splitext("/a/b.txt")                         // => "/a/b", ".txt"
ospath.Rel("/base", "/base/a.txt")                  // => "a.txt", nil
```

```go
file.Write("/tmp/a.txt", "x")                       // => nil
file.Read("/tmp/a.txt")                             // => "x", nil
file.Exists("/tmp/a.txt")                           // => true
file.ListDir("/tmp")                                // => []string, nil
file.Glob("/tmp/*.txt")                             // => []string, nil
```

## time / datetime

```go
t, _ := gotime.Strptime("2026-10-20 08:42:53", "%Y-%m-%d %H:%M:%S")
gotime.Strftime(t, "%Y年%m月%d日 %H:%M")             // => "2026年10月20日 08:42"
gotime.AddDays(t, 1)                                // => time.Time
gotime.DaysBetween(t, gotime.AddDays(t, 2))          // => 2
gotime.IsLeapYear(2024)                             // => true
gotime.Humanize(t)                                  // => "x 天前" 等（取决于当前时间）
```

```go
d, _ := datetime.NewDate(2026, time.October, 20)
d.ISOFormat()                                       // => "2026-10-20"
d.AddDays(1).ISOFormat()                            // => "2026-10-21"
datetime.NewTimedelta(36 * time.Hour).Days()        // => 1.5
datetime.FromTimestamp(0).Timestamp()               // => 0
```

## re

```go
re.Match(`\d+`, "123abc")                           // => true, nil
re.Search(`\d+`, "abc123")                          // => true, nil
re.FullMatch(`\d+`, "123")                          // => true, nil
re.FindAll(`\d+`, "a1b22c333")                      // => [1 22 333], nil
re.Sub(`(\w+)@(\w+)`, `${2}@${1}`, "tom@example")   // => "example@tom", nil
re.Split(`\s+`, "a b  c", 1)                        // => ["a" "b  c"], nil
re.FindStringSubmatchMap(`(?P<year>\d{4})-(?P<month>\d{2})`, "2024-06") // => map[month:06 year:2024], nil
```

## statistics / random

```go
statistics.Mean([]int{1, 2, 3})                     // => 2, nil
statistics.Median([]int{3, 1, 2})                   // => 2, nil
statistics.Mode([]int{1, 1, 2})                     // => 1, nil
statistics.Variance([]int{1, 2, 3})                 // => 1, nil
statistics.Stdev([]int{1, 2, 3})                    // => 1, nil
statistics.Quantiles([]int{1, 2, 3, 4}, 4)          // => [1.25 2.5 3.75], nil
statistics.Geometric_mean([]int{1, 4, 16})          // => 4, nil
statistics.Harmonic_mean([]int{1, 4, 4})            // => 2, nil
```

```go
random.Seed(42)
random.Float()                                      // => [0,1) 的 float64
random.Randint(1, 6)                                // => 1..6
random.Choice([]string{"a", "b"})                   // => "a" 或 "b", true
random.Sample([]int{1, 2, 3}, 2)                    // => []int, nil
random.Uniform(1, 2)                                // => [1,2] 的 float64
```

## 文本与编码

```go
textwrap.Wrap("The quick brown fox", 10)            // => ["The quick" "brown fox"]
textwrap.Fill("The quick brown fox", 10)            // => "The quick\nbrown fox"
textwrap.Dedent("\n    a\n      b\n")               // => "\na\n  b\n"
textwrap.Indent("a\n\nb", "> ")                     // => "> a\n\n> b"
textwrap.Shorten("The quick brown fox", 10)         // => "The qui..."
```

```go
base64.B64encode([]byte("hello"))                   // => "aGVsbG8="
base64.B64decode("aGVsbG8=")                        // => []byte("hello"), nil
html.Escape(`<b>"x"</b>`)                           // => "&lt;b&gt;&quot;x&quot;&lt;/b&gt;"
html.Unescape("&lt;b&gt;")                          // => "<b>"
```

```go
urllib.Quote("a b/c")                               // => "a%20b/c"
urllib.Unquote("a%20b")                             // => "a b", nil
urllib.QuotePlus("a b/c")                           // => "a+b%2Fc"
urllib.Encode(map[string][]string{"a": {"1", "2"}}) // => "a=1&a=2"
urllib.ParseQuery("a=1&a=2")                        // => map[a:[1 2]], nil
```

## string / textwrap

```go
string.Capwords("hello   world")                    // => "Hello World"
tmpl := string.NewTemplate("$name is ${age}")
tmpl.Substitute(map[string]string{"name": "Tom", "age": "18"}) // => "Tom is 18", nil
tmpl.Safe_substitute(map[string]string{"name": "Tom"})         // => "Tom is ${age}"
string.Digits                                       // => "0123456789"
```

## env / sys / logging

```go
env.Get("HOME")                                     // => 当前 HOME
env.Get("MISSING", "default")                       // => "default"
env.Set("APP_ENV", "dev")                           // => nil
env.Has("APP_ENV")                                  // => true
env.Expand("${APP_ENV}/x")                          // => "dev/x"
```

```go
sys.Args                                            // => []string，进程启动参数
sys.Getpid()                                        // => int
sys.Platform()                                      // => "darwin" / "linux" / "windows"
sys.Version()                                       // => "go1.x" 版本字符串
sys.Getcwd()                                        // => 当前目录, nil
```

```go
logger := logging.New("app")
logger.Info("started")                              // 输出: 2026-... [INFO] app started
logger.Warningf("slow: %dms", 120)                  // 输出: 2026-... [WARNING] app slow: 120ms
logging.BasicConfig(logging.LevelDebug)
logging.Debug("debug")
```

## subprocess / shutil / signal

```go
subprocess.CheckOutput("echo", "hello")             // => "hello\n", nil
subprocess.Run("echo", "hello").Stdout              // => "hello\n"
subprocess.RunShell("echo shell").Stdout            // => "shell\n"
```

```go
shutil.Copy("a.txt", "b.txt")                       // => nil
shutil.Copytree("src", "dst", true)                 // => nil
shutil.Rmtree("dst")                                // => nil
shutil.Which("go")                                  // => "/usr/local/go/bin/go", nil
shutil.DiskUsage("/")                               // => total, used, free, nil
```

```go
ch := make(chan os.Signal, 1)
signal.Notify(ch, signal.Interrupt)
signal.Stop(ch)

ctx, stop := signal.NotifyContext(context.Background(), signal.Interrupt)
stop()
```

## httpclient

```go
resp, err := httpclient.Get("https://example.com")
resp.StatusCode                                     // => 200
resp.OK()                                           // => true
resp.Text()                                         // => HTML 文本
err                                                 // => nil

httpclient.PostJSON("https://example.com/api", map[string]any{"name": "Tom"})
httpclient.GetJSON("https://example.com/data.json", &payload)
```

> HTTP 示例依赖网络和目标站点状态，返回码/文本仅供演示。

## conv / hash

```go
conv.Int("42")                                      // => 42, nil
conv.Float("1.5")                                   // => 1.5, nil
conv.Bool("true")                                   // => true, nil
conv.Bytes("abc")                                   // => []byte("abc"), nil
conv.ToJSON(map[string]any{"a": 1})                 // => "{\"a\":1}", nil
conv.FromJSON(`{"a":1}`, &out)                      // => nil
```

```go
hash.StringMD5("abc")                               // => "900150983cd24fb0d6963f7d28e17f72"
hash.StringSHA256("abc")                            // => "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
hash.UUID()                                         // => "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```
