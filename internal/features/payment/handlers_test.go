package payment

import (
	"testing"
	"time"
)

func TestResolvePaymentYear(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	years := []int{2027, 2026, 2025}

	t.Run("defaults to current year when available", func(t *testing.T) {
		got, err := resolvePaymentYear("", years, now)
		if err != nil {
			t.Fatalf("resolvePaymentYear() error = %v", err)
		}
		if got != 2026 {
			t.Fatalf("resolvePaymentYear() = %d, want 2026", got)
		}
	})

	t.Run("uses requested year", func(t *testing.T) {
		got, err := resolvePaymentYear("2025", years, now)
		if err != nil {
			t.Fatalf("resolvePaymentYear() error = %v", err)
		}
		if got != 2025 {
			t.Fatalf("resolvePaymentYear() = %d, want 2025", got)
		}
	})

	t.Run("falls back to latest when current year missing", func(t *testing.T) {
		got, err := resolvePaymentYear("", []int{2027, 2025}, now)
		if err != nil {
			t.Fatalf("resolvePaymentYear() error = %v", err)
		}
		if got != 2027 {
			t.Fatalf("resolvePaymentYear() = %d, want 2027", got)
		}
	})

	t.Run("invalid year returns error", func(t *testing.T) {
		_, err := resolvePaymentYear("abc", years, now)
		if err == nil {
			t.Fatal("resolvePaymentYear() expected error")
		}
	})
}
