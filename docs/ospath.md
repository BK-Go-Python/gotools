# ospath：路径字符串工具

`gotools/ospath` 提供 Python `os.path` 风格的工具。

```go
ospath.Join("a", "b", "c")
ospath.Base("/a/b.txt")   // b.txt
ospath.Dir("/a/b.txt")    // /a
ospath.Ext("/a/b.txt")    // .txt
root, ext := ospath.Splitext("/a/b.txt")
dir, file := ospath.Split("/a/b.txt")

ospath.Exists("a.txt")
ospath.IsAbs("/a")
ospath.Rel("/base", "/base/a.txt")
home, _ := ospath.Expanduser("~/x")
ospath.Expandvars("$HOME/x")
```

| Python | gotools |
|---|---|
| `os.path.join` | `Join` |
| `os.path.basename/dirname` | `Base` / `Dir` |
| `os.path.splitext` / `split` | `Splitext` / `Split` |
| `os.path.isabs/exists/isfile/isdir` | `IsAbs` / `Exists` / `IsFile` / `IsDir` |
| `os.path.abspath/relpath/normpath` | `Abs` / `Rel` / `Normpath` |
| `os.path.expanduser/expandvars` | `Expanduser` / `Expandvars` |
| `os.path.samefile` | `SameFile` |

## DIFF

- 底层是 `path/filepath`，路径分隔符跟随操作系统；Windows 使用 `\`。
- `Base`/`Dir` 遵循 Go `filepath` 语义，`Dir("a")` 返回 `"."`，Python 返回 `""`。
- `Expanduser` 只支持 `~` 与 `~/...`，不支持 `~user`。
- 没有 `commonpath`、`commonprefix`、`normcase`、`ismount`。
