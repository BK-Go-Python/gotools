# sys：脚本/系统接口

`gotools/sys` 提供 Python `sys`/`os` 常用脚本接口。

```go
fmt.Println(sys.Args)
fmt.Println(sys.Arg(1))
fmt.Println(sys.Getpid())
fmt.Println(sys.Platform())
fmt.Println(sys.Version())
fmt.Println(sys.Getcwd())
sys.Chdir("/tmp")
path, _ := sys.Executable()
host, _ := sys.Hostname()
```

| Python | gotools |
|---|---|
| `sys.argv` | `sys.Args` / `sys.Arg(i)` |
| `sys.exit(code)` | `sys.Exit` |
| `sys.platform` | `sys.Platform` |
| `sys.version` | `sys.Version` |
| `sys.maxsize` | `sys.Maxsize` |
| `os.getpid/getppid` | `sys.Getpid` / `sys.Getppid` |
| `os.getuid/geteuid/getgid/getegid` | `sys.Getuid` / `Geteuid` / `Getgid` / `Getegid` |
| `os.getcwd/chdir` | `sys.Getcwd` / `sys.Chdir` |
| `sys.executable` | `sys.Executable` |
| `socket.gethostname` | `sys.Hostname` |
| `sys.stdin/stdout/stderr` | `sys.Stdin/Stdout/Stderr` |

## DIFF

- `sys.Args` 是包初始化时从 `os.Args` 复制的，运行中修改 `os.Args` 不会自动同步。
- Windows 上 `Getuid/Geteuid/Getgid/Getegid` 返回 -1。
- `Exit` 会立即结束进程，测试中不要调用。
- `Version` 返回 Go 运行时版本，不是 CPython 版本字符串。
