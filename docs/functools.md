# functools：缓存工具

`gotools/functools` 提供 Python `functools` 中最常用的两个缓存工具。Go 没有装饰器语法，因此用高阶函数和泛型类型表达。

## Memoize

```go
fn := functools.Memoize(func(key int) int {
    return key * 2
})
fmt.Println(fn(2), fn(2)) // 4 4，第二次命中缓存
```

- 对应 Python `functools.cache`（无参数版本）
- 并发安全
- 同一个 key 的并发首次调用可能重复执行原函数，但最终只缓存一个结果

## LruCache

```go
cache := functools.NewLruCache[string, int](2)
cache.Put("a", 1)
cache.Put("b", 2)
cache.Get("a")      // 命中并把 a 移到最近使用
cache.Put("c", 3)   // 淘汰 b
cache.Keys()        // [c a]
```

- 对应 Python `functools.lru_cache(maxsize=capacity)`
- `Get` 命中会刷新 LRU 顺序；`Put` 超出容量淘汰最久未使用项
- `capacity <= 0` 时不存储
- 并发安全
- `Keys` 返回从最近使用到最久未使用的键

## Partial

```go
fn := func(a int, b string) string { return b + "!" }
partial, _ := functools.NewPartial(fn, 1)
result, err := partial.Call("x") // ["x!"]

add := func(a, b int) int { return a + b }
inc := functools.Partial2(add, 1) // func(int) int
fmt.Println(inc(2))               // 3
```

- `NewPartial(fn, bound...).Call(args...)`：反射版，返回 `[]any` 与 error
- `Partial1` / `Partial2`：类型安全的一元/二元偏函数
- 支持变参函数；panic 会转换为 error

## DIFF

- Go 版本没有装饰器，因此不提供 `@lru_cache` 语法；用 `Memoize`/`LruCache` 显式包装。
- `Memoize` 不统计命中率，也不暴露 `cache_info()`；需要统计时用 `LruCache`。
- `Partial.Call` 返回 `[]any`，不是直接返回强类型函数；强类型场景用 `Partial1`/`Partial2`。
- 不支持 Python 的 `partial.keywords`（关键字参数）；Go 没有关键字参数。
- `LruCache` 的容量是固定 int，不支持 `maxsize=None`（无界）；无界场景直接用 `Memoize` 或 map。
