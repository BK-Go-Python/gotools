// Package itertools 提供 Python itertools 风格的序列组合函数：
// 拼接、重复、循环、切片、条件截取、累积、相邻配对与最长拉链。
//
// 所有函数均以切片为输入并返回新切片，不修改入参；
// 结果为空时也返回非 nil 空切片。与 Python 的惰性迭代器不同，
// 本包函数会立即求出并返回完整结果。
//
// Python: import itertools
package itertools
