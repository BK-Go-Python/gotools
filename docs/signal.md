# signal：信号处理

`gotools/signal` 封装 Go `os/signal`，提供 Python `signal` 风格的入口。

```go
ch := make(chan os.Signal, 1)
signal.Notify(ch, signal.Interrupt)
go func() {
    <-ch
    fmt.Println("interrupted")
}()
// ...
signal.Stop(ch)
```

`NotifyContext` 适合优雅退出：

```go
ctx, stop := signal.NotifyContext(context.Background(), signal.Interrupt)
defer stop()

select {
case <-ctx.Done():
    // Ctrl-C
case <-workDone:
}
```

| Python | gotools |
|---|---|
| `signal.signal(signum, handler)` | `Notify(ch, sig)` + 读取 channel |
| `signal.getsignal`/`signal.signal` | 不直接对应 |
| `signal.default_int_handler` | `NotifyContext` |
| `signal.pause` | 阻塞读取 channel |
| `signal.SIG_IGN` | `Ignore` |
| `signal.SIG_DFL` | `Reset` |

## DIFF

- Go 使用 channel 而不是 handler 回调；Python 的 `signal.signal()` 返回旧 handler 的语义不适用。
- `NotifyContext` 同时监听父 context 取消；收到信号或父 context 取消都会 cancel。
- 信号处理是全局状态；测试中不要随意 `Ignore()`/`Reset()`，会影响整个进程。
- `Kill`/`Interrupt` 是 `os.Signal` 变量；其他信号用 `syscall.SIGTERM` 等。
- Windows 支持的信号集合与 Unix 不同。
