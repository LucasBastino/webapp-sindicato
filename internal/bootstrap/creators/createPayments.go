package creators

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/jmoiron/sqlx"
)

func CreatePayments(db *sqlx.DB) error {
	var pp []payment.Payment
	var p payment.Payment

	for i := range 50 {
		for j:= range 11 {
			p.Month = j+1
			p.Year = 2026
			p.Observations = ""
			p.CompanyID = i + 1
			p.DueDate = time.Date(p.Year, time.Month(p.Month), 15, 0, 0, 0, 0, time.UTC)
			random := rand.IntN(10)
			if random == 1 {
				p.PaidAt = nil
			} else {
				paidAt := time.Date(p.Year, time.Month(p.Month), 5, 0, 0, 0, 0, time.UTC)
				p.PaidAt = &paidAt
			}
			p.Amount = float32(rand.IntN(50000) + 50000)
		}
			pp = append(pp, p)
	}

	
	p.IsInPaymentPlan = false
	p.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))

	for _, p := range pp {
		insert, err := db.Query(`
		INSERT INTO payments(
		month,
		year,
		due_date,
		amount,
		paid_at,
		is_in_payment_plan,
		observations,
		id_company
		)
		VALUES (?,?,?,?,?,?,?,?)`,
			p.Month, p.Year, p.DueDate, p.Amount, p.PaidAt, p.IsInPaymentPlan, p.Observations, p.CompanyID)
		if err != nil {
			return fmt.Errorf("error inserting payments: %w", err)
		}
		insert.Close()
	}
	fmt.Println("payments created")
	return nil
}
