package creators

import (
	"testing"
	"time"
)

func TestIsFuturePaymentPeriod(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		year  int
		month int
		want  bool
	}{
		{name: "past month same year", year: 2026, month: 8, want: false},
		{name: "current month", year: 2026, month: 9, want: false},
		{name: "future month same year", year: 2026, month: 10, want: true},
		{name: "future year", year: 2027, month: 1, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFuturePaymentPeriod(tt.year, tt.month, now); got != tt.want {
				t.Fatalf("isFuturePaymentPeriod(%d, %d) = %v, want %v", tt.year, tt.month, got, tt.want)
			}
		})
	}
}

func TestIsFuturePaymentPeriod_May2027(t *testing.T) {
	now := time.Date(2027, time.May, 15, 0, 0, 0, 0, time.UTC)

	if isFuturePaymentPeriod(2027, 5, now) {
		t.Fatal("current month should not be future")
	}
	if !isFuturePaymentPeriod(2027, 6, now) {
		t.Fatal("june 2027 should be future when now is may 2027")
	}
}

func TestIsPendingPaymentPeriod(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		year  int
		month int
		want  bool
	}{
		{name: "past month", year: 2026, month: 8, want: false},
		{name: "current month", year: 2026, month: 9, want: true},
		{name: "future month", year: 2026, month: 10, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPendingPaymentPeriod(tt.year, tt.month, now); got != tt.want {
				t.Fatalf("isPendingPaymentPeriod(%d, %d) = %v, want %v", tt.year, tt.month, got, tt.want)
			}
		})
	}
}

func TestSeedPaymentPaidAt_FuturePeriodAlwaysNil(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	alwaysPaid := func(int) int { return 0 }

	cases := []struct {
		year  int
		month int
	}{
		{year: 2026, month: 10},
		{year: 2027, month: 1},
		{year: 2027, month: 12},
	}

	for _, tc := range cases {
		if paidAt := seedPaymentPaidAt(tc.year, tc.month, now, alwaysPaid); paidAt != nil {
			t.Fatalf("seedPaymentPaidAt(%d, %d) = %v, want nil", tc.year, tc.month, paidAt)
		}
	}
}

func TestSeedPaymentPaidAt_PastPeriodUsesRandom(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)

	if paidAt := seedPaymentPaidAt(2026, 9, now, func(int) int { return 0 }); paidAt != nil {
		t.Fatalf("current month should stay unpaid, got %v", paidAt)
	}

	alwaysUnpaid := func(int) int { return 1 }
	if paidAt := seedPaymentPaidAt(2026, 8, now, alwaysUnpaid); paidAt != nil {
		t.Fatalf("rng=1 should leave past month unpaid, got %v", paidAt)
	}

	alwaysPaid := func(int) int { return 0 }
	paidAt := seedPaymentPaidAt(2026, 8, now, alwaysPaid)
	if paidAt == nil {
		t.Fatal("rng=0 should mark past month as paid")
	}
	expected := time.Date(2026, time.August, 5, 0, 0, 0, 0, time.UTC)
	if !paidAt.Equal(expected) {
		t.Fatalf("paidAt = %v, want %v", paidAt, expected)
	}
}

func TestSeedPaymentPaidAt_May2027Scenario(t *testing.T) {
	now := time.Date(2027, time.May, 15, 0, 0, 0, 0, time.UTC)
	alwaysPaid := func(int) int { return 0 }

	if paidAt := seedPaymentPaidAt(2027, 6, now, alwaysPaid); paidAt != nil {
		t.Fatal("june 2027 should stay pending")
	}

	if paidAt := seedPaymentPaidAt(2027, 5, now, alwaysPaid); paidAt != nil {
		t.Fatal("may 2027 current month should stay pending")
	}

	paidAt := seedPaymentPaidAt(2027, 4, now, alwaysPaid)
	if paidAt == nil {
		t.Fatal("april 2027 should allow paid seed when rng marks paid")
	}
}

func TestSeedPaymentAmount_PendingPeriodIsNil(t *testing.T) {
	now := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	alwaysAmount := func(int) int { return 99999 }

	cases := []struct {
		year  int
		month int
	}{
		{year: 2026, month: 9},
		{year: 2026, month: 10},
		{year: 2027, month: 1},
	}

	for _, tc := range cases {
		if amount := seedPaymentAmount(tc.year, tc.month, now, alwaysAmount); amount != nil {
			t.Fatalf("seedPaymentAmount(%d, %d) = %v, want nil", tc.year, tc.month, amount)
		}
	}

	amount := seedPaymentAmount(2026, 8, now, alwaysAmount)
	if amount == nil {
		t.Fatal("past month should have amount")
	}
}
