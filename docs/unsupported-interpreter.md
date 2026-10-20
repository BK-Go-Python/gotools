# 解释器交互类内置函数：不提供 Go API

`compile`、`eval`、`exec`、`globals`、`locals`、`help`、`__import__` 依赖 CPython
解释器、模块系统和交互环境。Go 是静态编译语言，无法在同一进程内等价实现这些
语义。

本项目选择 **“不提供 API + 文档说明”**，而不是导出永远返回 `ErrNotSupported`
的空壳函数。原因有两点：

1. 项目约束禁止删除已有导出 API；空壳一旦加入就会永久存在。
2. 空壳没有真实行为，也没有可验证的测试价值，只会让调用方误以为“稍后可用”。

| Python 函数 | 典型用途 | 为什么 Go 无法模拟 | 建议替代 |
|---|---|---|---|
| `compile(source, filename, mode)` | 把源码编译成 code object | Go 没有运行时 Python 编译器；code object 是 CPython 内部结构 | 无等价物；动态逻辑用代码生成、配置 DSL 或独立 Python 进程 |
| `eval(expr)` | 求值表达式 | 同上，且需要 Python 作用域与对象模型 | 无等价物；不要把 `go/parser` 当作 Python 解析器 |
| `exec(code)` | 执行语句块 | 同上，还涉及模块/异常/命名空间 | 无等价物；需要脚本能力时用独立进程并明确部署依赖 |
| `globals()` | 获取模块全局命名空间 | Go 没有模块级可变命名空间 | 用显式 `struct`/`map` 作为上下文参数传递 |
| `locals()` | 获取局部命名空间 | Go 的局部变量不可反射枚举 | 用显式 `struct`/`map` 收集需要暴露的值 |
| `help(obj)` | 打开交互式帮助 | Go 没有 pydoc 运行时 | `go doc`、pkg.go.dev、IDE 文档 |
| `__import__("m")` | 运行时导入模块 | Go 的 import 是编译期静态的 | 使用 import 声明；动态加载用 `plugin`（平台受限）或独立进程 |

## 安全提示

- `eval`/`exec` 对不可信输入存在任意代码执行风险；本项目不提供这类入口。
- 不建议把 `os/exec` 调 Python 作为通用替代：它会引入 Python 运行时、依赖管理和
  跨平台部署成本，应该作为明确的架构选择而不是兼容层。

## 兼容性矩阵条目

> `compile`、`eval`、`exec`、`globals`、`locals`、`help`、`__import__`：不提供。
> Go 是静态编译语言，这些函数依赖 CPython 解释器/模块系统；请使用显式上下文、
> `go doc`/pkg.go.dev 和 Go 的 import 机制。详见 `docs/unsupported-interpreter.md`。
