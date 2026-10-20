# file：文件与路径操作

`gotools/file` 提供 Python `os`/`os.path`/`shutil` 风格的高频文件操作。

```go
file.Write("a.txt", "hello")
text, _ := file.Read("a.txt")
file.Copy("a.txt", "b.txt")
file.MkdirP("dir/sub")
list, _ := file.ListDir("dir")

info, _ := file.Stat("a.txt")
file.Chmod("a.txt", 0o600)
file.Walk("dir", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
})
```

| Python | gotools |
|---|---|
| `open(...).read()/write()` | `Read` / `Write` / `ReadBytes` / `WriteBytes` |
| `readlines` / `writelines` | `ReadLines` / `WriteLines` |
| `os.path.exists/is file/is dir` | `Exists` / `IsFile` / `IsDir` |
| `os.makedirs` | `MkdirP` |
| `os.listdir` / `glob.glob` | `ListDir` / `Glob` |
| `shutil.copy` / `shutil.move` | `Copy` / `Move` |
| `os.remove` / `shutil.rmtree` | `Remove` / `RemoveAll` |
| `os.stat` / `os.lstat` | `Stat` / `Lstat` |
| `os.chmod` / `os.chown` | `Chmod` / `Chown` |
| `os.rename` | `Rename` |
| `os.walk` | `Walk` / `WalkDir` |
| `os.path.islink` / `os.readlink` / `os.symlink` | `IsSymlink` / `Readlink` / `Symlink` |
| `os.path.samefile` | `SameFile` |
| `tempfile` | `TempFile` / `TempDir` |

## DIFF

- `Walk` 是回调式，不是 Python 的生成器 `(dirpath, dirnames, filenames)`；返回非 nil error 会终止遍历。
- `Chown` 在 Windows 上返回错误；普通用户通常只能改自己的文件。
- `Rename` 不跨文件系统移动；需要跨设备复制时用 `Move`。
- 没有 `shutil.copytree`、`rmtree` 的独立函数；递归删除用 `RemoveAll`。
