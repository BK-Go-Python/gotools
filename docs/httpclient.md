# httpclient：HTTP 请求

`gotools/httpclient` 提供 Python `urllib.request` 风格的 HTTP 客户端。

```go
resp, err := httpclient.Get("https://example.com/api")
fmt.Println(resp.StatusCode, resp.Text())

resp, err = httpclient.PostJSON("https://example.com/users", map[string]any{"name": "Tom"})

var data struct{ Value int }
err = httpclient.GetJSON("https://example.com/data", &data)
```

| Python | gotools |
|---|---|
| `urllib.request.urlopen(url)` | `httpclient.Get(url)` |
| `urllib.request.Request(url, data)` | `httpclient.Post` / `Request` |
| 读取响应体 | `Response.Body` / `Text()` |
| 解析 JSON | `Response.JSON(&v)` / `GetJSON` |
| `headers=` | `Options.Headers` |
| `timeout=` | `Options.Timeout` |
| 拼查询参数 | `Options.Query` |

## DIFF

- 非 2xx 状态码**不会**返回 error；用 `Response.OK()` 判断，Python `urlopen` 会抛 `HTTPError`。
- 默认超时 30 秒；`Options.Timeout <= 0` 时使用默认值。
- 没有重定向/代理/cookie jar 的显式配置；需要时直接用 `net/http`。
- 响应体一次性读入内存，不适合大文件下载。
- `PostJSON` 自动设置 `Content-Type: application/json`（如果调用方未设置）。
