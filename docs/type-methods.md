# Python 类型方法覆盖

`gotools` 覆盖的不只是 README 中按 Python 3.13 盘点的 71 个内置函数：`str`、`list`、`dict`、`set`/`frozenset` 的常用类型方法，以及 `int`、`float`、`bytes` 的常用方法，也都按 Python 同名 API 提供了 Go 实现。所有能力都以普通 Go 函数、泛型函数或具名类型的方法提供，例如 `str.Upper`、`(*list.List).Append`、`set.Set[T].Union`、`conv.BytesHex`；调用方直接 import 对应包即可，库本身不嵌入 Python 解释器，不链接 CPython，也不会通过 `os/exec` 调用外部 Python 来“模拟”兼容。

本文以 Python 3.13 的行为为基准。状态列的含义如下：

- **支持**：常用输入下的行为与 Python 一致，差异主要在 Go 的调用形式。
- **部分**：核心能力可用，但存在参数、返回类型、边界或语义差异。
- **不提供**：不导出空壳 API，需要改用 Go 标准库或其他方案。

DIFF 列只列最需要留意的差异；完整说明以各 API 的 GoDoc 中 `// DIFF:` 注释为准。示例省略 `import` 和错误处理，实际调用请检查 `error`。

## 总表

### str

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `upper()` | `str.Upper(s string) string` | 部分 | Go 使用简单大小写映射；仅额外补了 `ß → SS`，连字、希腊字母等完整 SpecialCasing 未覆盖 |
| `lower()` | `str.Lower(s string) string` | 部分 | 简单映射；仅补了 `İ → i + U+0307`，上下文相关规则不完整 |
| `capitalize()` | `str.Capitalize(s string) string` | 部分 | 首字符 titlecase、其余 lower；仅特判 `ß`，不是完整 case mapping |
| `title()` | `str.Title(s string) string` | 部分 | 用 `unicode.IsLetter` 判断词边界；titlecase 字母和特殊展开可能与 Python 不同 |
| `swapcase()` | `str.Swapcase(s string) string` | 部分 | 逐 rune 简单切换；仅特判 `ß`/`İ`，多字符映射缺失 |
| `casefold()` | `str.Casefold(s string) string` | 部分 | 仅特判 `ß`/`ẞ`/`İ`，其余退化为 `ToLower`；不是完整 Unicode case folding |
| `strip([chars])` | `str.Strip(s)`、`str.StripChars(s, chars)` | 支持 | 无参版按 Python 的空白集合；带 chars 的版本拆成独立函数 |
| `lstrip([chars])` | `str.Lstrip(s)`、`str.LstripChars(s, chars)` | 支持 | 同上 |
| `rstrip([chars])` | `str.Rstrip(s)`、`str.RstripChars(s, chars)` | 支持 | 同上 |
| `split(sep=None, maxsplit=-1)` | `str.Split(s, sep)`、`str.SplitN(s, sep, maxsplit)` | 部分 | `sep == ""` 时按空白分割；`maxsplit` 需改用 `SplitN`；空分隔符不会像 Python 那样抛 `ValueError` |
| `rsplit(sep=None, maxsplit=-1)` | `str.Rsplit(s, sep, maxsplit)` | 部分 | `maxsplit` 必填，传 `-1` 表示不限；`sep == ""` 时按空白分割 |
| `splitlines(keepends=False)` | `str.SplitLines(s)`、`str.SplitLinesKeepends(s)` | 支持 | 用两个函数区分 keepends；返回 `[]string` |
| `join(iterable)` | `str.Join(sep, []string)`、`str.JoinAny(sep, []interface{})` | 部分 | 只接受切片，不是任意可迭代；`JoinAny` 会先按 Python `str()` 转换元素，比 Python 宽松 |
| `find(sub[, start[, end]])` | `str.Find(s, sub) int` | 部分 | 返回 rune 下标，未找到返回 `-1`；不提供 start/end |
| `index(sub[, start[, end]])` | `str.Index(s, sub) (int, error)` | 部分 | 未找到返回 `error` 而非 `ValueError`；不提供 start/end |
| `count(sub[, start[, end]])` | `str.Count(s, sub) int` | 部分 | 统计不重叠出现；不提供 start/end |
| `startswith(prefix[, start[, end]])` | `str.Startswith(s string, prefixes ...string) bool` | 部分 | 变参模拟 Python 的前缀 tuple；不提供 start/end |
| `endswith(suffix[, start[, end]])` | `str.Endswith(s string, suffixes ...string) bool` | 部分 | 同上 |
| `center(width[, fillchar])` | `str.Center(s string, width int, fill ...rune) string` | 支持 | 按 rune 计宽；fill 省略时为空格 |
| `ljust(width[, fillchar])` | `str.Ljust(s string, width int, fill ...rune) string` | 支持 | 同上 |
| `rjust(width[, fillchar])` | `str.Rjust(s string, width int, fill ...rune) string` | 支持 | 同上 |
| `zfill(width)` | `str.Zfill(s string, width int) string` | 支持 | 按 rune 计宽；保留 ASCII `+`/`-` 的符号位置 |
| `replace(old, new[, count])` | `str.Replace(s, old, new string, count ...int) string` | 支持 | Go 变参；多传 count 只取第一个 |
| `translate(table)` | `str.Translate(s string, table map[rune]string) string` | 部分 | Python 用 None 删除；这里用空字符串表示删除 |
| `maketrans(from, to)` | `str.Maketrans(from, to string) (map[rune]string, error)` | 部分 | 仅支持 from/to 两个字符串；dict/三参数重载不覆盖 |
| `len(s)`、`s[i]`、`s[a:b]` | `str.Len`、`str.At`、`str.Slice` | 支持 | 按 rune 计数和索引；`Slice` 没有 step |
| `s[::-1]` | `str.Reverse(s)` | 支持 | 按 rune 反转，不按 grapheme cluster |
| `partition(sep)`、`rpartition(sep)` | `str.Partition`、`str.Rpartition` | 支持 | 返回 `(head, sep, tail)` 三个值；未找到时与 Python 一致 |
| `removeprefix(prefix)`、`removesuffix(suffix)` | `str.RemovePrefix`、`str.RemoveSuffix` | 支持 | 按字节前后缀匹配；合法 UTF-8 字符串下与 Python 结果一致 |
| `expandtabs(tabsize=8)` | `str.Expandtabs(s string, tabsize ...int)` | 支持 | 按 rune 列宽计算；`tabsize <= 0` 时删除制表符 |
| `is*()` 系列谓词 | `str.IsDigit`、`IsAlpha`、`IsSpace` 等 | 部分 | 使用 Go Unicode 表近似 Python 定义；`IsDigit` 不识别 `²` 等上标，`IsIdentifier` 非完整 XID |
| `str.format` 常用子集 | `str.Format`、`str.FormatMap`、`str.FormatValue` | 部分 | 只覆盖常用字段/格式说明符；不支持全部 Python 格式选项和 `__format__` |

### list

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `append(x)` | `(*list.List).Append(items ...interface{}) *list.List` | 部分 | 变参可一次追加多个；原地修改并返回接收者（Python 返回 `None`） |
| `extend(iterable)` | `(*list.List).Extend(v interface{}) error` | 部分 | 支持切片、数组、字符串和 map 键；字符串按 rune 展开，map 键排序；错误通过返回值而非 `TypeError` |
| `insert(i, x)` | `(*list.List).Insert(index int, value interface{}) *list.List` | 支持 | 负下标和越界截断规则与 Python 一致 |
| `remove(x)` | `(*list.List).Remove(value interface{}) error` | 部分 | 删除首个相等元素；不存在返回 `error` 而非 `ValueError` |
| `pop([i])` | `(*list.List).Pop()`、`PopAt(index int)` | 部分 | 拆成两个方法；空表或越界返回 `error`；`Pop()` 等价于 `pop(-1)` |
| `index(x[, start[, stop]])` | `List.Index(value)`、`List.IndexRange(value, start, stop)` | 部分 | 范围版必须显式传 start/stop；未找到返回 `error` 而非 `ValueError` |
| `count(x)` | `List.Count(value interface{}) int` | 支持 | 使用 Python 相等语义，包含 `1 == 1.0` 这类跨类型相等 |
| `sort(*, key=None, reverse=False)` | `(*List).Sort(less...)`、`SortBy(key, reverse...)`、`SortKey(key, reverse...)` | 部分 | 稳定排序，key 只计算一次；不可比较类型退化为 repr 比较而不是 `TypeError`；无关键字参数 |
| `reverse()` | `(*list.List).Reverse() *list.List` | 支持 | 原地反转 |
| `copy()` | `(*list.List).Copy() List` | 支持 | 浅拷贝 |
| `clear()` | `(*list.List).Clear() *list.List` | 支持 | 原地清空；保留底层数组容量，元素引用可能延迟释放 |
| `l[start:stop]` | `List.Slice(start, end) List` | 部分 | 没有 step，也不能省略参数；负数/越界按 Python 规则截断；返回新列表 |
| `len(l)`、`l[i]`、`l[i] = x` | `Len`、`Get`、`Set` | 支持 | 支持负下标 |
| `x in l`、`l1 == l2` | `Contains`、`Equal` | 支持 | 按 Python 语义做深比较 |
| `list(iterable)` | `list.New`、`list.From`、`list.MustFrom`、`list.Of[T]` | 部分 | 动态版返回 `List`；泛型版 `Of[T]` 返回同构的 `[]T` |

### dict

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `get(key[, default])` | `Dict.Get(key string, def ...interface{}) interface{}` | 部分 | 缺失时返回 `nil`/默认值，不抛 `KeyError`；默认值类型为 `interface{}` |
| `setdefault(key[, default])` | `Dict.SetDefault(key string, def interface{}) interface{}` | 支持 | 默认值必填，传 `nil` 等价于 Python 的 `None` |
| `keys()` | `Dict.Keys() []string` | 部分 | 按键排序返回快照；不是插入顺序，也不是动态视图 |
| `values()` | `Dict.Values() []interface{}` | 部分 | 顺序与排序后的 `Keys()` 一致 |
| `items()` | `Dict.Items() []Item` | 部分 | 顺序同上；`Item` 是含 `Key`/`Value` 字段的结构体，不是二元组 |
| `update(other)` | `Dict.Update(other interface{}) error` | 部分 | 只接受 mapping（`Dict`、`map[string]interface{}`、string 键 map）；不支持 pair 可迭代和关键字参数；错误通过返回值 |
| `pop(key[, default])` | `Dict.Pop(key string, def ...interface{}) interface{}` | 部分 | 缺失时返回 `nil`/默认值，不抛 `KeyError`；存在性请用 `Has` |
| `popitem()` | `Dict.PopItem() (Item, error)` | 部分 | 取排序后的第一个键，不是 Python 3.7+ 的 LIFO；空字典返回 `error` |
| `dict.fromkeys(keys[, value])` | `dict.FromKeys(keys []string, value ...interface{}) Dict` | 部分 | 包级函数而非类方法；键只能是 `[]string`；所有键共享同一个 value |
| `copy()` | `Dict.Copy() Dict` | 支持 | 浅拷贝 |
| `clear()` | `Dict.Clear() Dict` | 支持 | 原地清空 |
| 任意可哈希键 | — | 不提供 | `Dict` 固定为 `map[string]interface{}`；`dict.Of[K, V]` 只是泛型 map 浅拷贝，不是动态字典 |
| `d[k] = v`、`k in d`、`del d[k]` | `Set`、`Has`、`Delete` | 支持 | `Set` 返回 `Dict` 便于链式调用；`Delete` 返回原先是否存在 |
| `len(d)`、`d == other` | `Len`、`Equal` | 支持 | 深比较 |
| `dict(mapping)` | `dict.New`、`dict.From`、`dict.Of[K, V]` | 部分 | 动态版只支持 string 键 |
| 嵌套路径读取（扩展） | `Dict.GetPath` | 部分 | 非 Python 标准方法；支持点号路径与切片下标，仅作为便利 API |

### set / frozenset

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `set.add(x)` | `Set[T].Add(items ...T) Set[T]` | 部分 | 变参；对零值 `Set` 调用 `Add` 会返回新 map，必须接收返回值 |
| `set.discard(x)` | `Set[T].Discard(items ...T) Set[T]` | 支持 | 缺失元素忽略；历史方法 `Remove` 保留相同的 discard 语义 |
| `set.remove(x)` | `Set[T].RemoveOne(item T) error` | 部分 | 缺失返回 `error`（对应 `KeyError`）；命名为 `RemoveOne`，避免与 discard 语义的 `Remove` 冲突 |
| `set.union(*others)` | `Set[T].Union(other Set[T]) Set[T]` | 部分 | 只接受一个 other；返回新集合 |
| `set.intersection(*others)` | `Set[T].Intersection(other Set[T]) Set[T]` | 部分 | 只接受一个 other |
| `set.difference(*others)` | `Set[T].Difference(other Set[T]) Set[T]` | 部分 | 只接受一个 other |
| `set.symmetric_difference(other)` | `Set[T].SymmetricDifference(other Set[T]) Set[T]` | 部分 | 只接受一个 other；不接受任意 iterable |
| `set.issubset(other)` | `Set[T].IsSubset(other Set[T]) bool` | 支持 | 等价于 `<=`；不提供真子集 `<` |
| `set.issuperset(other)` | `Set[T].Issuperset(other Set[T]) bool` | 支持 | 等价于 `>=`；注意方法名是 `Issuperset` |
| `set.isdisjoint(other)` | `Set[T].IsDisjoint(other Set[T]) bool` | 支持 | 无 |
| `set.update(*others)` | `Set[T].Update(others ...Set[T]) Set[T]` | 部分 | 原地并入并返回集合；只接受 `Set`，不接受任意 iterable |
| `set.intersection_update(*others)` | `Set[T].IntersectionUpdate(other Set[T]) Set[T]` | 部分 | 只接受一个 other |
| `set.difference_update(*others)` | `Set[T].DifferenceUpdate(other Set[T]) Set[T]` | 部分 | 只接受一个 other |
| `set.symmetric_difference_update(other)` | `Set[T].SymmetricDifferenceUpdate(other Set[T]) Set[T]` | 支持 | 只接受一个 other |
| `set.pop()` | `Set[T].Pop() (T, bool)` | 部分 | 移除并返回任意元素；空集合返回零值和 `false`，不抛 `KeyError` |
| `set.clear()` | `Set[T].Clear() Set[T]` | 支持 | 原地清空 |
| `set.copy()` | `Set[T].Copy() Set[T]`、`Clone()` | 支持 | 浅拷贝 |
| `frozenset` | `Frozen[T]`、`set.NewFrozen` | 部分 | 不可变集合；没有 add/remove/update/pop/clear；`Copy`/`Clone` 返回自身 |
| `frozenset` 的非变更运算 | `Frozen[T].Union`、`Intersection`、`Difference`、`SymmetricDifference`、`IsSubset`、`Issuperset`、`IsDisjoint` | 支持 | 与 `Set` 同名同义，返回 `Frozen[T]`；没有原地更新方法 |
| `set(iterable)`、`frozenset(iterable)` | `set.New`、`set.NewFrozen` | 部分 | 只接受变参，从切片构造需要显式循环 |
| 元素类型要求 | `T comparable` | 部分 | 不能放切片、map、函数；Python 只要求元素 hashable（可用 tuple/frozenset） |
| 集合运算符（并、交、差、对称差、子集、超集） | — | 不提供 | Go 不支持运算符重载；请使用 `Union`、`Intersection`、`Difference`、`SymmetricDifference`、`IsSubset`、`Issuperset` |

### int / float

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `int.bit_length()` | `py.BitLength[T Integer](x T) int` | 支持 | 支持定宽 int/uint 系列；负数按绝对值计算；不支持任意精度整数 |
| `float.is_integer()` | `conv.IsInteger(f float64) bool` | 支持 | NaN/Inf 返回 `false`；函数形式而非方法；没有 Decimal 版本 |
| `int.to_bytes(length, byteorder, *, signed=False)` | `conv.IntToBytes(x int64, length int, order string, signed ...bool)` | 部分 | int64；length 1 到 1 MiB；Python 支持任意精度/长度 |
| `int.from_bytes(bytes, byteorder, *, signed=False)` | `conv.IntFromBytes(data []byte, order string, signed ...bool)` | 部分 | 仅 1–8 字节、返回 int64；超出 int64 报错 |
| `int.bit_count()` | — | 不提供 | 可改用 `math/bits.OnesCount64`；gotools 未封装 |

### bytes

| Python 类型/方法 | gotools API | 状态 | 主要 DIFF |
|---|---|---|---|
| `bytes.hex([sep[, bytes_per_sep]])` | `conv.BytesHex(data []byte) string` | 部分 | 返回小写十六进制字符串；不提供 sep/bytes_per_sep；返回 string 而非 bytes |
| `bytes.fromhex(s)` | `conv.FromHex(s string) ([]byte, error)` | 部分 | 只忽略空格、Tab、CR、LF，未忽略全部 ASCII 空白；奇数长度或非法字符返回 `error`；类方法用包函数替代 |
| `bytes(v)`、`bytearray(v)` | `conv.Bytes`、`conv.ByteArray` | 部分 | 支持 string、`[]byte`、整数和字节切片；不提供 encoding 参数 |
| `bytes.decode()`、`str.encode()` | — | 不提供 | 直接用 `string(b)` / `[]byte(s)`；其他编码转换需要 `golang.org/x/text` 等外部依赖 |

## 三个重点差异

### 1. 有序性：Go map 无序，`dict` 输出按键排序

Python 3.7+ 的 `dict` 保留插入顺序，`popitem()` 默认按 LIFO 弹出；Go 的 map 不保留任何顺序，迭代顺序是随机的。`gotools` 的 `dict.Dict` 就是 `map[string]interface{}`，因此**无法复现插入顺序**，只能选择一种确定性的替代方案：

- `Keys()` 用 `sort.Strings` 对键排序后返回；
- `Values()`、`Items()` 复用排序后的键，因此三者的顺序始终一致；
- `PopItem()` 返回排序后的**第一个**键值对，而不是 Python 的最后一个；
- `String()`/Python 风格 repr 同样按键排序，便于生成稳定的测试输出；
- `list.From(map)` 得到的键切片也按键排序，避免 map 随机顺序影响结果。

```go
d := dict.From(map[string]interface{}{"b": 2, "a": 1})
fmt.Println(d.Keys()) // [a b]：排序结果，不是插入顺序

item, _ := d.PopItem()
fmt.Println(item.Key, item.Value) // a 1：固定取排序后的第一个
```

`set.Items()` 直接遍历底层 map，顺序不保证；这一点与 Python `set` 的“顺序任意”类似，但需要稳定输出时应显式排序：

```go
items := set.New(3, 1, 2).Items()
slices.Sort(items)
fmt.Println(items) // [1 2 3]
```

如果业务依赖插入顺序（例如按字段定义顺序序列化、需要 LIFO 栈语义），不要用 `dict.Dict` 承担顺序：应额外维护 `[]string` 键顺序，或直接用有序的结构体切片/键值对切片。

### 2. 整数精度：Go 定宽溢出 vs Python 任意精度

Python 的 `int` 是任意精度整数，运算结果会自动扩展；Go 的 `int`、`int64`、`uint64` 是定宽整数，普通算术溢出会静默回绕。`gotools` 没有用 `math/big` 替换所有整数 API，因此需要区分两类情况：

- **转换路径会报错**：`conv.Int`、`conv.Int64` 对超出目标范围的字符串/浮点数返回 `error`，不会截断；`conv.Uint64` 拒绝负数。
- **算术路径可能溢出**：`py.Abs` 对最小整数取负会溢出（GoDoc 已注明）；`py.Sum` 统一以 `float64` 累加，大整数会丢精度；`py.SumOf[T]` 在 `T` 中累加，溢出时按 Go 规则回绕；`py.Pow` 使用 `math.Pow`，本身是浮点运算。
- **位运算只支持定宽类型**：`py.BitLength` 覆盖 int/uint 系列，不接受 `big.Int`；Python 的 `int.bit_length()` 可用于任意大整数。
- **字节转换同样受限**：`conv.IntToBytes`/`IntFromBytes` 以 int64 和 1–8 字节为边界，不能像 Python 那样处理任意精度、任意长度的整数。

```go
_, err := conv.Int64("9223372036854775808") // 超过 MaxInt64
fmt.Println(err != nil) // true

fmt.Println(py.BitLength(-8)) // 4
```

需要任意精度时，直接使用 `math/big`：把大整数运算留在 `big.Int` 中，只在确定范围后才转回 `int64`。金额等场景优先使用最小货币单位（如分）的整数或定点 decimal，不要依赖浮点累加。对可能溢出的累加/乘法，应显式使用 `big.Int` 或自行做边界检查。

### 3. Unicode：Go 简单大小写映射 vs Python 完整 case mapping/folding

Go 的 `string` 是 UTF-8 字节序列，Python 的 `str` 是 Unicode 码点序列。`gotools` 在涉及“字符”的 API 上统一按 rune（码点）处理：`str.Len`、`str.At`、`str.Slice`、`str.Find`、`str.Reverse` 都以 rune 为单位，因此下标语义与 Python 的码点下标一致；`str.LenBytes` 则有意返回 UTF-8 字节数。两者都不按 grapheme cluster（用户感知字符）切分，组合字符、ZWJ emoji 仍会被拆成多个码点。

大小写是更大的差异。Go 的 `unicode.ToUpper/ToLower/ToTitle` 是**简单大小写映射**，通常一个码点映射到另一个码点；Python 使用完整的 Unicode SpecialCasing 和 case folding，包含：

- 一对多展开，如 `ß → SS`、`ﬄ → FFL`、连字和部分希腊字母；
- 上下文相关规则，如希腊词尾 sigma；
- 完整的 case folding，如 `ẞ → ss`、`ς → σ`、连字折叠。

`gotools` 在 `Upper`、`Lower`、`Capitalize`、`Title`、`Swapcase`、`Casefold` 中补了 `ß`、`ẞ`、`İ` 等常见特例，但没有完整复制 SpecialCasing 表，因此仍属于“部分”支持。`IsIdentifier`、`IsPrintable`、`IsDigit` 等谓词也使用 Go Unicode 表的近似实现，不等同于 Python 的完整定义。

```go
fmt.Println(str.Len("世界"))       // 2（rune/码点）
fmt.Println(str.LenBytes("世界"))  // 6（UTF-8 字节）

fmt.Println(str.Upper("ß"))        // SS（已特判）
fmt.Println(str.Upper("ﬄ"))        // ﬄ：完整映射未覆盖
```

此外，`str.Center`、`Ljust`、`Rjust`、`Zfill` 按码点数而不是终端显示宽度计算，这与 Python 的 `len` 语义一致，但东亚宽字符在等宽终端中仍会显得更宽。

实践建议：对 ASCII 或已知字符集，`gotools` 的大小写方法足够；需要与 Python 完全一致的 Unicode case mapping/folding 时，应限制输入字符集，或引入 `golang.org/x/text`（会打破本项目的零第三方依赖），并显式做 Unicode 规范化。`gotools` 不做 NFC/NFD 规范化，包含组合字符的字符串应按码点而不是“看起来是否相同”来比较。

## 使用建议

| 场景 | 建议 |
|---|---|
| 从 Python 迁移脚本、需要 Python 风格语义 | 使用 `gotools`：真值判断、repr、向下取整除法、rune 下标、排序后的 dict 输出等都已封装 |
| 数据类型异构、原型开发、配置/JSON 边界处理 | 使用 `list.List`、`dict.Dict` 和 `conv`；注意它们以 `interface{}` 为载体，会带来装箱和运行时断言成本 |
| 生成稳定的文本输出、测试快照 | 使用 `dict.Keys/Items`、`list.Sort` 的稳定排序；`set.Items` 需要显式排序 |
| 热路径、强类型数据处理 | 使用 Go 标准库：`slices.SortFunc`、`slices.Contains`、`maps.Keys`、`strings.Builder`，避免 `interface{}` 装箱 |
| 对顺序敏感的逻辑 | 不要依赖 `dict`；额外维护键顺序，或使用结构体切片/有序 map |
| 大整数、金额、累加可能溢出 | 使用 `math/big` 或定点整数；不要用 `py.Sum` 做精确整数求和 |
| 完整 Unicode 大小写/规范化 | 限制为 ASCII，或引入 `golang.org/x/text`；不要假设 `str.Casefold` 等价于 Python `casefold()` |
| 大文件/流式处理 | 使用 `bufio`、`io.Reader`、`iter.Seq`；`gotools` 的许多 API 会先物化为切片 |

混用时的几条原则：

1. **在边界转换，不在循环里来回转换。** `list.List` 是 `[]interface{}`，`list.Of[T]` 是 `[]T`；如果内部逻辑是强类型切片，不要为了调用一个 `gotools` 方法而反复装箱/拆箱。
2. **把 `error` 当作 Python 异常处理。** `conv.Int64`、`list.Remove`、`dict.Update`、`conv.FromHex` 等都通过返回值报错；不要忽略，否则会把失败当成零值。
3. **注意原地修改语义。** `list` 的 `Append/Insert/Sort/Reverse` 等修改接收者并返回指针；`set.Add` 对零值会返回新 map，必须写回；`dict` 的方法直接修改底层 map。
4. **需要 Python 关键字参数时，改用分离的方法。** 例如 `Sort`/`SortBy`、`Strip`/`StripChars`、`Pop`/`PopAt`、`IsSubset`/`Issuperset`；Go 没有关键字参数和运算符重载。
5. **并发写入需要外部同步。** `list.List`、`dict.Dict`、`set.Set` 都不是并发安全容器；Go map 的并发写还会直接触发运行时错误。
6. **需要插入顺序时，显式建模。** 用 `[]string` 保存键顺序，或使用 `[]struct{Key, Value}`；`dict` 的排序输出只是确定性方案，不是 Python 插入顺序的替代品。

一个简短的组合示例：

```go
d := dict.From(map[string]interface{}{"b": 2, "a": 1})
fmt.Println(d.Keys()) // [a b]

l := list.New(3, 1, 2)
l.SortBy(func(v interface{}) interface{} { return v }, true)
fmt.Println(l) // [3 2 1]

s := set.New(1, 2, 3).Intersection(set.New(3, 4))
items := s.Items()
slices.Sort(items)
fmt.Println(items) // [3]

b, _ := conv.FromHex("de ad be ef")
fmt.Println(conv.BytesHex(b)) // deadbeef
```

## 延伸阅读

- [`docs/unsupported-interpreter.md`](unsupported-interpreter.md)：`eval`、`exec`、`globals` 等依赖 CPython 解释器的函数为何不提供。
- [`docs/unsupported-protocol.md`](unsupported-protocol.md)：异步迭代、描述符、`property`/`super` 等对象协议的边界与替代方案。
