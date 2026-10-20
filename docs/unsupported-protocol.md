# 异步、描述符与对象协议：部分实现或不提供

| Python 函数 | 状态 | 说明 | 建议替代 |
|---|---|---|---|
| `breakpoint()` | 部分实现 | `py.Breakpoint()` 调用 `runtime.Breakpoint()`，触发底层陷阱；不会启动 pdb | 使用 Delve 等 Go 调试器 |
| `memoryview(b)` | 部分实现 | `py.MemoryView([]byte)` 返回共享底层数组的一维视图 | 需要多维/格式时使用 `[]byte` + 显式偏移 |
| `aiter(obj)` | 不提供 | Go 没有 async 迭代协议 | 用 `iter.Seq`/channel 表达流 |
| `anext(obj)` | 不提供 | 同上 | 用 `for range` 或 `<-chan` |
| `classmethod` | 不提供 | Go 没有类方法与描述符协议 | 用普通函数或包级函数 |
| `staticmethod` | 不提供 | 同上 | 用普通函数 |
| `property` | 不提供 | Go 没有属性描述符 | 用 getter/setter 方法或导出字段 |
| `super()` | 不提供 | Go 没有继承模型 | 用接口组合与显式嵌入 |
| `object()` | 不提供 | Go 没有统一的可变对象/属性模型 | 用 `struct{}{}` 或具体结构体 |

## 设计说明

- **不导出空壳**：`aiter`、`anext`、`classmethod`、`property`、`staticmethod`、
  `super`、`object` 不会出现在 `py` 包中。它们无法给出诚实行为，加入空壳会违反
  “不删除现有导出 API”的约束。
- **`breakpoint` 不写测试**：调用 `runtime.Breakpoint()` 会让测试进程陷入 trap，
  因此只提供函数和文档。
- **`memoryview` 只做视图**：不复制数据，修改视图会影响原 `[]byte`；`nil` 返回
  `nil`，空切片返回空切片。

## 兼容性矩阵条目

> `breakpoint`：部分实现（`py.Breakpoint`，调用 `runtime.Breakpoint`）。
> `memoryview`：部分实现（`py.MemoryView`，仅一维 `[]byte` 视图）。
> `aiter`、`anext`、`classmethod`、`property`、`staticmethod`、`super`、`object`：不提供。
> 详见 `docs/unsupported-protocol.md`。
