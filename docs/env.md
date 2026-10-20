# env：环境变量工具

`gotools/env` 提供 Python `os.environ` 风格的环境变量操作。

```go
env.Set("APP_ENV", "prod")
fmt.Println(env.Get("APP_ENV"))          // prod
fmt.Println(env.Get("MISSING", "dev"))   // dev
fmt.Println(env.Has("APP_ENV"))          // true
fmt.Println(env.Expand("${APP_ENV}/x"))  // prod/x
env.Unset("APP_ENV")
```

| Python | gotools |
|---|---|
| `os.environ.get(k[, default])` | `env.Get(k[, default])` |
| `os.environ[k]` / `k in os.environ` | `env.Lookup` / `env.Has` |
| `os.environ[k] = v` | `env.Set` |
| `del os.environ[k]` | `env.Unset` |
| `os.environ.keys()` | `env.Keys`（排序） |
| `os.environ.items()` | `env.All`（`KEY=VALUE`，排序） |
| `os.path.expandvars` | `env.Expand`（`$K` 与 `${K}`） |

## DIFF

- `All` 返回排序后的 `KEY=VALUE` 切片，不是插入顺序的映射视图。
- `Expand` 对未定义变量保留原样（Go `os.ExpandEnv` 行为）；Python 的 `expandvars` 也保留。
- 环境变量在进程内是全局可变状态，测试/并发代码需自行同步。
