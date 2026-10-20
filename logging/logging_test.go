package logging

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := New("app")
	logger.SetOutput(&buf)
	logger.SetLevel(LevelWarning)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warning("warn")
	logger.Errorf("error %d", 7)

	output := buf.String()
	if strings.Contains(output, "debug") || strings.Contains(output, "info") {
		t.Fatalf("low level logs leaked: %q", output)
	}
	if !strings.Contains(output, "[WARNING] app warn") || !strings.Contains(output, "[ERROR] app error 7") {
		t.Fatalf("output = %q", output)
	}
}

func TestPackageLogger(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer func() {
		SetOutput(nil)
		SetLevel(LevelInfo)
	}()

	Infof("hello %s", "world")
	Critical("boom")
	output := buf.String()
	if !strings.Contains(output, "[INFO] hello world") || !strings.Contains(output, "[CRITICAL] boom") {
		t.Fatalf("output = %q", output)
	}
}

func TestLoggerConcurrent(t *testing.T) {
	var buf bytes.Buffer
	logger := New("c")
	logger.SetOutput(&buf)
	logger.SetLevel(LevelDebug)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info("x")
		}()
	}
	wg.Wait()
	if got := strings.Count(buf.String(), "[INFO] c x"); got != 20 {
		t.Fatalf("concurrent logs = %d", got)
	}
}
