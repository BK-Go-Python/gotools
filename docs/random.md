# random：Python random 风格随机数

`gotools/random` 包装 `math/rand`，提供常用随机函数。全局生成器由互斥锁保护，并发安全。

```go
random.Seed(42)
fmt.Println(random.Float())         // [0,1)
fmt.Println(random.Randint(1, 6))   // 1..6
fmt.Println(random.Randrange(0, 10, 2))
fmt.Println(random.Choice([]string{"a", "b"}))
items := []int{1, 2, 3}
random.Shuffle(items)
fmt.Println(random.Sample(items, 2))
```

| Python | gotools |
|---|---|
| `random.seed(x)` | `random.Seed` |
| `random.random()` | `random.Float` |
| `random.uniform(a,b)` | `random.Uniform` |
| `random.randint(a,b)` | `random.Randint` |
| `random.randrange(start,stop,step)` | `random.Randrange` |
| `random.choice(seq)` | `random.Choice` |
| `random.shuffle(seq)` | `random.Shuffle` |
| `random.sample(seq,k)` | `random.Sample` |
| `random.choices(seq,k=k)` | `random.Choices` |

## DIFF

- `Randint(a,b)` 在 `a > b` 时自动交换；Python 抛 ValueError。
- `Choice` 空输入返回 `(零值,false)`；Python 抛 IndexError。
- `Sample`/`Choices` 非法参数返回 error；Python 抛 ValueError。
- 不提供 `gauss`、`betavariate`、`seed` 的版本参数等分布函数。
- 随机算法是 Go `math/rand`，不保证与 CPython 的 Mersenne Twister 序列一致；只保证同一次运行内 `Seed` 可复现。
