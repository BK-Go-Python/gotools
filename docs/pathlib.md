# pathlib：Path 对象

`gotools/pathlib` 提供 `Path` 具名类型，覆盖 Python `pathlib.Path` 的高频操作。

```go
p := pathlib.New("data", "users.json")
text, err := p.ReadText()
err = p.WriteText(`{"ok":true}`)
exists := p.Exists()
items, err := p.Parent().Iterdir()
matches, err := p.Parent().Glob("*.json")
```

| Python | gotools |
|---|---|
| `Path(*parts)` | `pathlib.New(parts...)`、`pathlib.From(s)` |
| `p / "a"` | `p.Join("a")` |
| `p.parent` / `p.name` / `p.suffix` / `p.stem` | `Parent` / `Name` / `Suffix` / `Stem` |
| `p.with_suffix(".md")` | `WithSuffix` |
| `p.exists/is_file/is_dir` | `Exists` / `IsFile` / `IsDir` |
| `p.read_text/write_text` | `ReadText` / `WriteText` |
| `p.read_bytes/write_bytes` | `ReadBytes` / `WriteBytes` |
| `p.mkdir(parents=True)` | `Mkdir(true)` |
| `p.iterdir()` / `p.glob("*.go")` | `Iterdir` / `Glob` |
| `p.stat().st_size` | `Size` |
| `p.touch()` / `p.rename()` / `p.unlink()` / `p.rmdir()` | `Touch` / `Rename` / `Unlink` / `Rmdir` |
| `p.resolve()` | `Abs` |

## DIFF

- Go 没有运算符重载，`p / "a"` 改为 `p.Join("a")`。
- `Mkdir` 用变参 `parents ...bool` 表达 `parents=`；`exist_ok` 未提供，已存在会返回 error。
- `Glob` 底层是 `filepath.Glob`，语法与 Python `pathlib` 基本一致，但 `**` 递归语义不同。
- `ReadText/WriteText` 固定 UTF-8；其他编码需要显式转换。
- `Path` 是 string 的具名类型，不做路径规范化；`New` 使用 `filepath.Join`。
