# html：转义与反转义

`gotools/html` 提供 Python `html` 模块的 `escape`/`unescape`。

```go
html.Escape(`<a href="x">'&`)
// &lt;a href=&quot;x&quot;&gt;&#x27;&amp;

html.Escape(`<b>`, false)
// &lt;b&gt;

html.Unescape("&lt;b&gt;")
// <b>
```

## DIFF

- `Escape` 默认 `quote=true`，与 Python 一致；单引号转义为 `&#x27;`（Go `html.EscapeString` 使用 `&#39;`）。
- `Unescape` 直接使用 Go `html.UnescapeString`，支持常见命名实体和数字实体。
