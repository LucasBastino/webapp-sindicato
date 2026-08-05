package payment

import (
	"testing"
	"time"
)

func TestGetStatus(t *testing.T){
	now := time.Now()
	paidAt := now.Add(-24*time.Hour)
	cases := []struct{
		name string
		p Payment
		want string
	}{
		{name: "valid completed at time", p: Payment{DueDate:now.Add(48*time.Hour), PaidAt: &paidAt, IsInPaymentPlan: false}, want: "Completado"},
		{name: "valid completed in overdue", p: Payment{DueDate:now.Add(-48*time.Hour), PaidAt: &paidAt, IsInPaymentPlan: false}, want: "Completado"},
		{name: "valid overdue", p: Payment{DueDate:now.Add(-48*time.Hour), PaidAt: nil, IsInPaymentPlan: false}, want: "Vencido"},
		{name: "valid pending", p: Payment{DueDate:now.Add(48*time.Hour), PaidAt: nil, IsInPaymentPlan: false}, want: "Pendiente"},
		{name: "valid in payment plan", p: Payment{DueDate:now.Add(-48*time.Hour), PaidAt: &paidAt, IsInPaymentPlan: true}, want: "En plan de pago"},
	}

	for _, c := range cases{
		t.Run(c.name, func(t *testing.T) {
			got := c.p.GetStatus()
			if got != c.want{
				t.Fatalf("GetStatus() = %q, want %q", got, c.want)
			}
		})
	}
}