# shutil：文件树工具

`gotools/shutil` 提供 Python `shutil` 风格的文件/目录树操作。

```go
shutil.Copy("a.txt", "b.txt")
shutil.Copy2("a.txt", "dir/")
shutil.Copytree("src", "dst")
shutil.Copytree("src", "dst", true) // dirs_exist_ok
shutil.Move("dst", "moved")
shutil.Rmtree("moved")

path, err := shutil.Which("go")
total, used, free, err := shutil.DiskUsage("/")
```

| Python | gotools |
|---|---|
| `shutil.copy` | `Copy` |
| `shutil.copy2` | `Copy2` |
| `shutil.copytree` | `Copytree(src, dst, dirsExistOK...)` |
| `shutil.rmtree` | `Rmtree(path, ignoreErrors...)` |
| `shutil.move` | `Move` |
| `shutil.which` | `Which` |
| `shutil.disk_usage` | `DiskUsage`（Windows 未实现） |

## DIFF

- `Copy`/`Copy2` 的 `dst` 是目录时复制到目录下，与 Python 一致。
- `Copytree` 不复制符号链接的目标，按普通文件/目录处理；没有 `symlinks`、`ignore` 参数。
- `Move` 优先 `os.Rename`，跨设备时复制后删除。
- `DiskUsage` 在 Windows 返回 error。
- 没有 `make_archive`/`unpack_archive`、`copymode`/`copystat` 独立函数。
