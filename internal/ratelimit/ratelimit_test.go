package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNew_NoLimit(t *testing.T) {
	lim := New(0, 0)
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("expected no error for disabled limiter, got %v", err)
	}
	if !lim.Allow() {
		t.Fatal("expected Allow() to return true for disabled limiter")
	}
}

func TestNew_Allowed(t *testing.T) {
	lim := New(10, 5)
	defer func() {
		if err := lim.Wait(context.Background()); err != nil {
			t.Fatalf("unexpected error on cleanup: %v", err)
		}
	}()
	if !lim.Allow() {
		t.Fatal("expected Allow() to return true")
	}
}

func TestNew_WaitBlocks(t *testing.T) {
	lim := New(100, 1)
	// Consume the burst token.
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Next Wait should block since burst is 1 and refill rate is 100/s.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- lim.Wait(ctx)
	}()
	// Wait a bit to confirm goroutine is blocked.
	select {
	case <-done:
		t.Fatal("expected Wait to block")
	case <-time.After(50 * time.Millisecond):
		// Still blocked, as expected.
	}
	// Cancel the context to unblock.
	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("goroutine did not exit after context cancel")
	}
}

func TestNew_WaitAfterRefill(t *testing.T) {
	lim := New(10, 1)
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Wait for the token to refill (10 per second = 100ms per token).
	time.Sleep(120 * time.Millisecond)
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error after refill: %v", err)
	}
}

func TestNew_WaitTimeout(t *testing.T) {
	lim := New(100, 1)
	// Consume the burst token.
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Use a short timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := lim.Wait(ctx)
	if err == nil {
		t.Fatal("expected error from timed-out context")
	}
}

func TestReserve(t *testing.T) {
	lim := New(10, 1)
	res := lim.Reserve()
	if res == nil {
		t.Fatal("expected non-nil Reservation")
	}
	defer res.Cancel()
}

func TestSetLimit(t *testing.T) {
	lim := New(10, 5)
	lim.SetLimit(20)
	lim.SetBurst(10)
	// Should not panic.
	_ = lim.Allow()
}

func TestAllow_AfterWait(t *testing.T) {
	lim := New(100, 2)
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := lim.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Burst exhausted, Allow may be false after tokens are consumed.
	_ = lim.Allow()
}