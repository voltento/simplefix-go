package simplefixgo

import (
	"context"
	"errors"
	"testing"
	"time"
)

// A taker that drains slower than the acceptor produces reports fills the
// outbound channel and never lets it empty. sendRaw must not park the sender
// past the configured deadline: it returns a timeout and cancels the handler
// context, the same escape a fully stalled socket already gets.
func TestSendRawTimesOutWhenOutboundStaysFull(t *testing.T) {
	const deadline = 50 * time.Millisecond

	h := NewAcceptorHandler(context.Background(), "35", 1)
	h.SetSendDeadline(deadline)

	if err := h.SendRaw([]byte("first")); err != nil {
		t.Fatalf("filling the outbound channel: %v", err)
	}

	started := time.Now()
	err := h.SendRaw([]byte("second"))
	elapsed := time.Since(started)

	if !errors.Is(err, ErrSendTimeout) {
		t.Fatalf("sendRaw error = %v, want %v", err, ErrSendTimeout)
	}
	if elapsed < deadline || elapsed > 5*deadline {
		t.Fatalf("sendRaw returned after %s, want about %s", elapsed, deadline)
	}

	select {
	case <-h.Context().Done():
	case <-time.After(time.Second):
		t.Fatal("handler context was not canceled on send timeout")
	}
}

func TestSendRawHonoursZeroDeadline(t *testing.T) {
	h := NewAcceptorHandler(context.Background(), "35", 0)

	done := make(chan error, 1)
	go func() { done <- h.SendRaw([]byte("never drained")) }()

	select {
	case err := <-done:
		t.Fatalf("sendRaw returned %v with no deadline set, want it to wait", err)
	case <-time.After(100 * time.Millisecond):
	}

	h.Stop()
	if err := <-done; err == nil {
		t.Fatal("sendRaw returned nil after the handler was stopped")
	}
}
