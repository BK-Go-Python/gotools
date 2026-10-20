# itertools：有限迭代工具

`gotools/itertools` 把 Python `itertools` 的常用函数实现为返回切片的泛型函数。Go 没有无限惰性迭代器，因此多数函数需要显式长度或返回已物化的结果。

| Python | gotools | 说明 |
|---|---|---|
| `chain(*seqs)` | `Chain(seqs...)` | 依次拼接 |
| `repeat(value, n)` | `Repeat(value, n)` | 重复 n 次 |
| `cycle(items)` + `islice` | `Cycle(items, n)` | 循环取 n 个 |
| `islice(items, start, stop, step)` | `Islice` | 非负 start/stop、正 step |
| `takewhile` / `dropwhile` | `Takewhile` / `Dropwhile` | 前缀处理 |
| `accumulate` | `Accumulate` | 需要显式累积函数 |
| `pairwise` | `Pairwise` | 相邻二元组 |
| `zip_longest` | `Zip_longest` | 缺失位置填 fill |
| `product` | `Product` | 笛卡尔积 |
| `combinations` / `permutations` | `Combinations` / `Permutations` | 按索引字典序 |
| `combinations_with_replacement` | `Combinations_with_replacement` | 允许重复取元素 |
| `groupby` | `Groupby` | 仅连续分组 |
| `filterfalse` | `Filterfalse` | 保留 pred=false |
| `compress` | `Compress` | 布尔选择器 |
| `starmap` | `Starmap` | 使用 `Pair2` |
| `count` | `Count(start, step, n)` | 有限版本 |
| `tee` | `Tee(items, n)` | 返回 n 份独立副本 |

## DIFF

- Python 返回惰性迭代器；gotools 返回切片，需要预先知道长度。
- `Cycle`、`Count`、`Tee` 是有限版本，必须传入个数/长度。
- `Groupby` 只合并相邻的相同 key，与 Python 一致；不是全局分组。
- `Starmap` 使用 `itertools.Pair2`，不是任意长度元组。
