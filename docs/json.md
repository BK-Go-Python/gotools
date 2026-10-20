# json：Python json 风格 API

`gotools/json` 包装 `encoding/json`，提供 Python 常用的 `loads/dumps/load/dump/indent`。

```go
value, err := json.Loads(`{"a":1}`)
text, err := json.Dumps(map[string]any{"a": 1})       // {"a":1}
pretty, err := json.DumpsIndent(value, "  ")
err = json.Dump("data.json", value)
value, err = json.Load("data.json")
```

| Python | gotools |
|---|---|
| `json.loads(s)` | `json.Loads(s)` |
| `json.load(fp)` | `json.Load(path)` |
| `json.dumps(v)` | `json.Dumps(v)` |
| `json.dumps(v, indent=2)` | `json.DumpsIndent(v, "  ")` |
| `json.dump(v, fp)` | `json.Dump(path, v)` |
| `json.loads` + `dumps` | `json.Indent(s, "  ")` |
| 解析到结构体 | `json.LoadsInto(s, &out)` |

## DIFF

- 整数优先转 `int64`；超出 `int64` 的整数保留为 `json.Number`（Python 是任意精度）。
- 默认不开启 `ensure_ascii`，中文原样输出。
- 不支持 `object_hook`、`parse_float`、`default` 等 Python 参数；需要自定义时用 `encoding/json`。
- `loads` 会拒绝尾随数据，与 Python 一致。
