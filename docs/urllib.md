# urllib：URL 编解码与解析

`gotools/urllib` 提供 Python `urllib.parse` 的高频函数。

```go
urllib.Quote("a b/c")          // a%20b/c
urllib.Quote("a b/c", "")      // a%20b%2Fc
urllib.Unquote("a%20b")        // a b

urllib.QuotePlus("a b/c")      // a+b%2Fc
urllib.UnquotePlus("a+b")      // a b

urllib.Encode(map[string][]string{"a": {"1", "x y"}})
urllib.ParseQuery("a=1&a=2")
urllib.Parse("https://example.com/a?x=1")
urllib.Join("https://example.com/a/", "b")
```

| Python | gotools |
|---|---|
| `urllib.parse.quote` | `Quote(s, safe...)` |
| `urllib.parse.unquote` | `Unquote` |
| `urllib.parse.quote_plus` | `QuotePlus` |
| `urllib.parse.unquote_plus` | `UnquotePlus` |
| `urllib.parse.urlencode` | `Encode` |
| `urllib.parse.parse_qs` | `ParseQuery` |
| `urllib.parse.urlparse` | `Parse`（返回 `*net/url.URL`） |
| `urllib.parse.urljoin` | `Join` |

## DIFF

- `Quote` 默认 safe 是 `/`；`QuotePlus` 默认 safe 是空字符串，与 Python 一致。
- `Unquote` 对非法 `%` 转义返回 error；Python 保留原样。
- `Parse` 返回 Go `*url.URL`，字段名与 Python `ParseResult` 不同。
- `Join` 使用 `url.JoinPath`，路径语义与 Python `urljoin` 在 `..`、查询串等边界上可能不同。
- 没有 `urlsplit/urlunsplit/parse_qsl/urlunparse` 等完整 API。
