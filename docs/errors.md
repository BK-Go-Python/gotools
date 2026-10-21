# 错误与 panic

Go 没有 Python 的 `try/except`，只有 panic/recover。gotools 的策略是：

1. **新 API 优先返回 error**，不 panic。
2. **历史动态 API 保持兼容**：它们可能 panic；需要安全调用时用 `E` 变体或 `py.Try`。
3. **不静默吞错**：不会为了“统一捕获”把 panic 变成零值继续执行。

## 通用捕获

```go
value, err := py.Try(func() []int {
    return py.Range(1, 2, 3, 4) // 参数个数非法，会 panic
})
fmt.Println(value, err) // [] <nil 或 recovered panic>

err = py.Catch(func() {
    py.Len(42) // Len 对非法类型 panic
})
```

- `py.Try[T]` / `py.Try2[A,B]` / `py.Try3[A,B,C]`：捕获 panic、返回 error 和零值
- `py.Catch(fn)`：只捕获 panic、返回 error

## 安全变体

| 旧 API | 安全 API | 说明 |
|---|---|---|
| `py.Range` | `py.RangeE` | 参数个数/step 非法返回 error |
| `py.FloorDiv` / `py.DivMod` | `py.FloorDivE` / `py.DivModE` | 除零返回 `ErrZeroDivision` |
| `py.Mod` | `py.ModE` | 除零返回 `ErrZeroDivision` |
| `py.Reduce` | `py.ReduceE` | 空序列无初始值返回 `ErrEmpty`；回调 panic 转 error |
| `py.Any` / `py.All` | `py.AnyE` / `py.AllE` | 非法 iterable 返回 error |
| `py.EnumerateAny` | `py.EnumerateAnyE` | 非法 iterable 返回 error |
| `py.Zip` | `py.ZipE` | 非法 iterable 返回 error |
| `py.SortedAny` / `py.SortedBy` | `py.SortedAnyE` / `py.SortedByE` | 非法 iterable 或比较回调 panic 转 error |
| `py.ReversedAny` | `py.ReversedAnyE` | 非法 iterable 返回 error |
| `py.MapAny` / `py.FilterAny` | `py.MapAnyE` / `py.FilterAnyE` | 非法 iterable 或回调 panic 转 error |
| `py.Len` | `py.LenE` | 非法类型返回 error |

`SumAny`、`MinAny`、`MaxAny` 本身返回 error，现在对非法 iterable 也不再 panic。

## 为什么不给所有方法加 recover

- 大量方法没有 `error` 返回值（如 `str.Title`、`list.Append`）。强行 recover 只能返回零值，会掩盖真实 bug。
- 给 600+ 个 API 改签名会破坏兼容性，也让正常调用变得冗长。
- 正确做法是：输入校验用 error；对旧动态 API 用 `E` 变体；第三方/未知回调用 `Try`/`Catch`。

`py.Len`、`py.Range`、`py.Reduce`、`py.FloorDiv`、`py.Mod` 的旧版本仍保留 panic 行为，避免已有代码语义变化。
