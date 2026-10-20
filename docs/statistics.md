# statistics：常用统计函数

`gotools/statistics` 提供 Python `statistics` 模块的高频函数，接受任意数值切片。

```go
mean, _ := statistics.Mean([]int{1, 2, 3})       // 2
median, _ := statistics.Median([]int{3, 1, 2})   // 2
mode, _ := statistics.Mode([]int{1, 1, 2})       // 1
modes, _ := statistics.Multimode([]int{1, 1, 2, 2}) // [1 2]
variance, _ := statistics.Variance([]int{1, 2, 3})  // 1
stdev, _ := statistics.Stdev([]int{1, 2, 3})        // 1
```

| Python | gotools |
|---|---|
| `statistics.mean` | `Mean` |
| `statistics.median` | `Median` |
| `statistics.mode` | `Mode` |
| `statistics.multimode` | `Multimode` |
| `statistics.variance` | `Variance`（样本，n-1） |
| `statistics.pvariance` | `Pvariance`（总体，n） |
| `statistics.stdev` | `Stdev` |
| `statistics.pstdev` | `Pstdev` |

## DIFF

- 返回值统一为 `float64`；Python 对部分输入返回 `Fraction`。
- `Mode` 并列时按 `fmt.Sprint` 排序取第一个；Python 3.8+ 返回首个遇到的值。
- 空输入、样本方差只有一个数据点时返回 error；Python 抛 `StatisticsError`。
- 不提供 `quantiles`、`fmean`、`geometric_mean`、`harmonic_mean`、`covariance` 等长尾函数。
