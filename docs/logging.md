# logging：轻量日志

`gotools/logging` 提供 Python `logging` 风格的日志记录器。输出格式：

```text
2026-10-20 16:30:00 [INFO] app started
```

```go
logger := logging.New("app")
logger.SetLevel(logging.LevelDebug)
logger.Info("started")
logger.Warningf("slow: %dms", 120)

logging.BasicConfig(logging.LevelInfo)
logging.Info("hello")
```

## API

- `New(prefix)`：创建 `*Logger`，默认 stderr、Info
- `SetLevel` / `SetOutput` / `SetPrefix`
- `Debug/Info/Warning/Warn/Error/Critical`
- `Debugf/Infof/Warningf/Errorf/Criticalf`
- `Log/Logf`
- 包级 `SetLevel/SetOutput/BasicConfig` 与同名日志函数

## 级别

`LevelDebug`、`LevelInfo`、`LevelWarning`、`LevelError`、`LevelCritical`；
同时提供 Python 拼写别名 `DEBUG`、`INFO`、`WARNING`、`ERROR`、`CRITICAL`。

## DIFF

- 没有 Python 的 handler/formatter/filter/logger 层级注册表；用 `Logger` 实例和 `SetOutput` 组合。
- 没有 `logging.getLogger(name)` 全局单例；`New(prefix)` 返回独立记录器。
- 时间格式固定为 `2006-01-02 15:04:05`，前缀为空时不输出前缀。
- 级别常量用 `LevelInfo` 等，避免与包级 `Info()` 函数重名；`INFO` 等别名可用。
