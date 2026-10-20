# re：Go RE2 风格正则

`gotools/re` 提供 Python `re` 模块的高频入口，但底层是 Go `regexp`（RE2）。

```go
ok, err := re.Match(`\d+`, "123abc")       // true，锚定开头
ok, err = re.Search(`\d+`, "abc123")       // true
all, err := re.FindAll(`\d+`, "a1b22")     // [1 22]
out, err := re.Sub(`(\w+)\s+(\w+)`, `${2} ${1}`, "hello world")
parts, err := re.Split(`\s+`, "a b  c", 1)
```

| Python | gotools |
|---|---|
| `re.match` | `re.Match`（只匹配开头） |
| `re.fullmatch` | `re.FullMatch` |
| `re.search` | `re.Search` |
| `re.findall`（无分组） | `re.FindAll`（只返回整体匹配） |
| `re.findall`（有分组） | `re.FindAllStringSubmatch` |
| `match.group(0..n)` | `re.FindStringSubmatch` / `Regexp.FindStringSubmatch` |
| `re.search(...).group(0)` | `re.FindString` |
| `re.sub` | `re.Sub` |
| `re.split` | `re.Split` |
| `re.compile` | `re.Compile`，返回 `*Regexp` |
| `re.escape` | `re.QuoteMeta` |

## DIFF

- Go RE2 不支持回溯引用、lookahead/lookbehind、条件表达式；依赖这些语法的 Python 正则需要重写。
- 替换模板使用 Go 的 `$1`、`${name}`，不是 Python 的 `\1`。
- `FindAll` 只返回整体匹配；需要分组时用 `FindAllStringSubmatch`（下标 0 是整体，1..n 是分组）。
- `re.Match`/`Search` 返回 bool；需要匹配文本时用 `FindString` 或 `Compile` 后的方法。
- Python 的 `re.MULTILINE`/`DOTALL` 等标志用 `(?m)`/`(?s)` 内联表达。
