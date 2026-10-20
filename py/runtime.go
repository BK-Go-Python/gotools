package py

import "runtime"

// Breakpoint 等价于 Python 的 breakpoint()，触发底层断点陷阱。
//
// Python: breakpoint()
// DIFF: 调用 runtime.Breakpoint，不会启动 pdb 或任何交互式调试器；
// 是否中断取决于调试器/信号处理配置。
func Breakpoint() {
	runtime.Breakpoint()
}
