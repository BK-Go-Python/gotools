// Package signal 提供 Python signal 模块风格的 Go 信号封装。
//
// Go 的信号模型基于 channel 而不是 Python 的 handler 回调；Notify 把信号
// 转发到 channel，NotifyContext 在收到信号或父 context 取消时取消子 context。
package signal

import (
	"context"
	"os"
	stdsignal "os/signal"
)

// Signal 是 os.Signal 的别名。
type Signal = os.Signal

// 常用信号值。
var (
	Interrupt Signal = os.Interrupt
	Kill      Signal = os.Kill
)

// Notify 等价于 Go 的 signal.Notify，把 sigs 转发到 c。
func Notify(c chan<- os.Signal, sigs ...os.Signal) {
	stdsignal.Notify(c, sigs...)
}

// Stop 停止向 c 转发信号。
func Stop(c chan<- os.Signal) {
	stdsignal.Stop(c)
}

// Ignore 忽略 sigs；不传参数时忽略所有可忽略信号。
func Ignore(sigs ...os.Signal) {
	stdsignal.Ignore(sigs...)
}

// Reset 撤销 Ignore/Notify 对 sigs 的影响；不传参数时重置所有信号。
func Reset(sigs ...os.Signal) {
	stdsignal.Reset(sigs...)
}

// NotifyContext 等价于 Go 的 signal.NotifyContext：收到信号或父 context
// 取消时取消返回的 context。
func NotifyContext(parent context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	return stdsignal.NotifyContext(parent, sigs...)
}
