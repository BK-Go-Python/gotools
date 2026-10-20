# textwrap：文本换行与缩进

`gotools/textwrap` 提供 Python `textwrap` 的高频函数。

```go
textwrap.Wrap("The quick brown fox", 10)
// ["The quick" "brown fox"]

textwrap.Fill("The quick brown fox", 10)
// "The quick\nbrown fox"

textwrap.Dedent("\n    a\n      b\n    c\n")
textwrap.Indent("a\n\nb", "> ")
textwrap.Shorten("The quick brown fox", 10) // "The qui..."
```

| Python | gotools |
|---|---|
| `textwrap.wrap` | `Wrap` |
| `textwrap.fill` | `Fill` |
| `textwrap.dedent` | `Dedent` |
| `textwrap.indent` | `Indent`（只给非空白行加前缀） |
| `textwrap.shorten` | `Shorten` |

## DIFF

- `Wrap` 先按空白切词，不保留原有换行/缩进；Python 默认也会替换空白。
- 超长单词直接按宽度断开（Python 的 `break_long_words=True`）。
- `width <= 0` 返回空切片；Python 抛 `ValueError`。
- 不提供 `TextWrapper` 类、`initial_indent`/`subsequent_indent`/`max_lines` 等完整选项。
