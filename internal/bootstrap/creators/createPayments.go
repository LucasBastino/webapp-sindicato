package creators

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/jmoiron/sqlx"
)

func CreatePayments(db *sqlx.DB) error {
	var pp []payment.Payment

	for i := range 50 {
		for j := range 12 {
			var p payment.Payment
			p.Month = j + 1
			p.Year = 2027
			p.IsInPaymentPlan = false
			p.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))
			p.CompanyID = i + 1
			dueDate := time.Date(p.Year, time.Month(p.Month), 15, 0, 0, 0, 0, time.UTC)
			p.DueDate = dueDate
			random := rand.IntN(10)
			if random == 1 {
				p.PaidAt = nil
			} else {
				paidAt := time.Date(p.Year, time.Month(p.Month), 5, 0, 0, 0, 0, time.UTC)
				p.PaidAt = &paidAt
			}
			amt := float32(rand.IntN(50000) + 50000)
			p.Amount = &amt
			pp = append(pp, p)
		}
	}

	if len(pp) == 0 {
		return nil
	}

	query := `INSERT INTO payments(
		month,
		year,
		due_date,
		amount,
		paid_at,
		is_in_payment_plan,
		observations,
		id_company
	) VALUES`
	placeholders := make([]string, 0, len(pp))
	args := make([]any, 0, len(pp)*8)
	for _, p := range pp {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?)")
		args = append(args, p.Month, p.Year, p.DueDate, p.Amount, p.PaidAt, p.IsInPaymentPlan, p.Observations, p.CompanyID)
	}
	query += strings.Join(placeholders, ",")

	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error inserting payments: %w", err)
	}
	fmt.Println("payments created")
	return nil
}

// func CreatePayments(db *sqlx.DB) error {
// 	var pp []payment.Payment
// 	var p payment.Payment
//
// 	for i := range 50 {
// 		for j:= range 12 {
// 			p.Month = j+1
// 			p.Year = 2026
// 			p.IsInPaymentPlan = false
// 			p.Observations = fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999))
// 			p.CompanyID = i + 1
// 			dueDate := time.Date(p.Year, time.Month(p.Month), 15, 0, 0, 0, 0, time.UTC)
// 			p.DueDate = &dueDate
// 			random := rand.IntN(10)
// 			if random == 1 {
// 				p.PaidAt = nil
// 			} else {
// 				paidAt := time.Date(p.Year, time.Month(p.Month), 5, 0, 0, 0, 0, time.UTC)
// 				p.PaidAt = &paidAt
// 			}
// 			amt := float32(rand.IntN(50000) + 50000)
// 			p.Amount = &amt
// 			pp = append(pp, p)
// 		}
// 	}
//
//
//
// 	for _, p := range pp {
// 		insert, err := db.Query(`
// 		INSERT INTO payments(
// 		month,
// 		year,
// 		due_date,
// 		amount,
// 		paid_at,
// 		is_in_payment_plan,
// 		observations,
// 		id_company
// 		)
// 		VALUES (?,?,?,?,?,?,?,?)`,
// 			p.Month, p.Year, p.DueDate, p.Amount, p.PaidAt, p.IsInPaymentPlan, p.Observations, p.CompanyID)
// 		if err != nil {
// 			return fmt.Errorf("error inserting payments: %w", err)
// 		}
// 		insert.Close()
// 	}
// 	fmt.Println("payments created")
// 	return nil
// }
