package signal

import (
	"context"
	"os"
	"testing"
)

func TestNotifyContext(t *testing.T) {
	ctx, cancel := NotifyContext(context.Background(), Interrupt)
	if ctx == nil || cancel == nil {
		t.Fatal("NotifyContext returned nil")
	}
	cancel()
	select {
	case <-ctx.Done():
	default:
		t.Fatal("context was not canceled")
	}
}

func TestNotifyStop(t *testing.T) {
	ch := make(chan os.Signal, 1)
	Notify(ch, Interrupt)
	Stop(ch)

	var _ Signal = Interrupt
	var _ Signal = Kill
}
