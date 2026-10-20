// Package collections 提供 Python collections 模块风格的容器。
//
// 本包使用 Go 的类型系统和标准库实现 Python 同名容器的常用语义，
// 不嵌入 Python 解释器，也不依赖任何第三方库。Counter 是
// collections.Counter 的泛型版本，可以直接 range、len 和下标访问。
//
// Go 的 map 不保留插入顺序，因此凡是 Python 依赖插入顺序或“任意顺序”
// 的 API（例如 Counter.MostCommon 的并列项、Counter.Elements），这里都会
// 按稳定规则排序；具体差异见对应方法上的 // DIFF 注释。
package collections
