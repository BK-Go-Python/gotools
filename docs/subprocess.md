# subprocess：命令执行

`gotools/subprocess` 封装 `os/exec`，提供 Python `subprocess` 风格的高频入口。

```go
result, err := subprocess.Run("echo", "hello")
fmt.Println(result.Stdout) // hello\n
fmt.Println(result.Returncode)

out, err := subprocess.CheckOutput("sh", "-c", "printf hi")
err = subprocess.CheckCall("sh", "-c", "exit 0")
result, err = subprocess.RunShell("echo $HOME")
```

| Python | gotools |
|---|---|
| `subprocess.run(args)` | `subprocess.Run(name, args...)` |
| `subprocess.run(..., capture_output=True)` | 默认捕获 stdout/stderr 到 `Result` |
| `subprocess.check_output(args)` | `subprocess.CheckOutput` |
| `subprocess.check_call(args)` | `subprocess.CheckCall` |
| `subprocess.run(..., shell=True)` | `subprocess.RunShell` |
| `timeout=` / `cwd=` / `env=` / `input=` | `Options.Timeout` / `Dir` / `Env` / `Stdin` |

## DIFF

- 默认不经过 shell；需要 shell 时显式用 `RunShell`，仅用于可信输入。
- `Run` 对非零退出码返回 `Result` 和 `nil` error（与 `check=False` 一致）；`CheckOutput/CheckCall` 才在非零时返回 error。
- `Options.Env` 非 nil 时**替换**整个环境，不会自动继承；需要继承时传 `append(os.Environ(), "K=V")`。
- 没有 `Popen` 流式管道、`communicate`、`stdin=PIPE` 等完整进程控制；长驻进程请直接用 `os/exec`。
- 超时通过 `context.WithTimeout` 实现，超时后子进程会被 Kill。
