// Package logging 提供 Python logging 风格的轻量日志封装。
//
// 输出格式为：日期 时间 [LEVEL] prefix message；默认级别为 Info。
package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level 是日志级别，数值越大越严重。
type Level int

const (
	LevelDebug    Level = 10
	LevelInfo     Level = 20
	LevelWarning  Level = 30
	LevelError    Level = 40
	LevelCritical Level = 50
)

// 兼容 Python 常量拼写。
const (
	DEBUG    = LevelDebug
	INFO     = LevelInfo
	WARNING  = LevelWarning
	ERROR    = LevelError
	CRITICAL = LevelCritical
)

// String 返回级别名称。
func (l Level) String() string {
	switch {
	case l >= LevelCritical:
		return "CRITICAL"
	case l >= LevelError:
		return "ERROR"
	case l >= LevelWarning:
		return "WARNING"
	case l >= LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

// Logger 是日志记录器，并发安全。
type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	level  Level
	prefix string
}

// New 创建日志记录器，默认输出到 stderr、级别 Info。
func New(prefix string) *Logger {
	return &Logger{out: os.Stderr, level: LevelInfo, prefix: prefix}
}

var std = New("")

// SetLevel 设置日志级别。
func (l *Logger) SetLevel(level Level) { l.level = level }

// SetOutput 设置输出目标。
func (l *Logger) SetOutput(w io.Writer) { l.out = w }

// SetPrefix 设置日志前缀。
func (l *Logger) SetPrefix(prefix string) { l.prefix = prefix }

// Log 写入一条指定级别的日志。
func (l *Logger) Log(level Level, args ...any) {
	l.log(level, fmt.Sprint(args...))
}

// Logf 写入格式化日志。
func (l *Logger) Logf(level Level, format string, args ...any) {
	l.log(level, fmt.Sprintf(format, args...))
}

// Debug 写入 DEBUG 日志。
func (l *Logger) Debug(args ...any) { l.Log(LevelDebug, args...) }

// Info 写入 INFO 日志。
func (l *Logger) Info(args ...any) { l.Log(LevelInfo, args...) }

// Warning 写入 WARNING 日志。
func (l *Logger) Warning(args ...any) { l.Log(LevelWarning, args...) }

// Warn 是 Warning 的别名。
func (l *Logger) Warn(args ...any) { l.Warning(args...) }

// Error 写入 ERROR 日志。
func (l *Logger) Error(args ...any) { l.Log(LevelError, args...) }

// Critical 写入 CRITICAL 日志。
func (l *Logger) Critical(args ...any) { l.Log(LevelCritical, args...) }

// Debugf 写入格式化 DEBUG 日志。
func (l *Logger) Debugf(format string, args ...any) { l.Logf(LevelDebug, format, args...) }

// Infof 写入格式化 INFO 日志。
func (l *Logger) Infof(format string, args ...any) { l.Logf(LevelInfo, format, args...) }

// Warningf 写入格式化 WARNING 日志。
func (l *Logger) Warningf(format string, args ...any) { l.Logf(LevelWarning, format, args...) }

// Errorf 写入格式化 ERROR 日志。
func (l *Logger) Errorf(format string, args ...any) { l.Logf(LevelError, format, args...) }

// Criticalf 写入格式化 CRITICAL 日志。
func (l *Logger) Criticalf(format string, args ...any) { l.Logf(LevelCritical, format, args...) }

func (l *Logger) log(level Level, message string) {
	if level < l.level {
		return
	}
	out := l.out
	if out == nil {
		out = os.Stderr
	}
	prefix := ""
	if l.prefix != "" {
		prefix = " " + l.prefix
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(out, "%s [%s]%s %s\n",
		time.Now().Format("2006-01-02 15:04:05"), level.String(), prefix, message)
}

// SetLevel 设置默认日志级别。
func SetLevel(level Level) { std.SetLevel(level) }

// SetOutput 设置默认输出目标。
func SetOutput(w io.Writer) { std.SetOutput(w) }

// BasicConfig 配置默认日志器；output 省略时保持当前输出。
func BasicConfig(level Level, output ...io.Writer) {
	std.SetLevel(level)
	if len(output) > 0 {
		std.SetOutput(output[0])
	}
}

// Debug 写入默认日志器的 DEBUG 日志。
func Debug(args ...any) { std.Debug(args...) }

// Info 写入默认日志器的 INFO 日志。
func Info(args ...any) { std.Info(args...) }

// Warning 写入默认日志器的 WARNING 日志。
func Warning(args ...any) { std.Warning(args...) }

// Warn 是 Warning 的别名。
func Warn(args ...any) { std.Warn(args...) }

// Error 写入默认日志器的 ERROR 日志。
func Error(args ...any) { std.Error(args...) }

// Critical 写入默认日志器的 CRITICAL 日志。
func Critical(args ...any) { std.Critical(args...) }

// Debugf 写入默认日志器的格式化 DEBUG 日志。
func Debugf(format string, args ...any) { std.Debugf(format, args...) }

// Infof 写入默认日志器的格式化 INFO 日志。
func Infof(format string, args ...any) { std.Infof(format, args...) }

// Warningf 写入默认日志器的格式化 WARNING 日志。
func Warningf(format string, args ...any) { std.Warningf(format, args...) }

// Errorf 写入默认日志器的格式化 ERROR 日志。
func Errorf(format string, args ...any) { std.Errorf(format, args...) }

// Criticalf 写入默认日志器的格式化 CRITICAL 日志。
func Criticalf(format string, args ...any) { std.Criticalf(format, args...) }
