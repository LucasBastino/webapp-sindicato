package creators

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	"github.com/jmoiron/sqlx"
)

func CreatePayments(db *sqlx.DB) error {
	companyIDs, err := listActiveCompanyIDs(db, 50)
	if err != nil {
		return err
	}
	if len(companyIDs) == 0 {
		return nil
	}

	now := time.Now()
	currentYear := now.Year()
	years := []int{currentYear -1, currentYear, currentYear + 1}

	var pp []payment.Payment
	for _, companyID := range companyIDs {
		for _, year := range years {
			for month := 1; month <= 12; month++ {
				p := payment.Payment{
					Month:           month,
					Year:            year,
					IsInPaymentPlan: false,
					Observations:    fmt.Sprintf("texto aleatorio numero %d", rand.IntN(999)),
					CompanyID:       companyID,
					DueDate:         time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC),
				}
				p.PaidAt = seedPaymentPaidAt(year, month, now, rand.IntN)
				p.Amount = seedPaymentAmount(year, month, now, rand.IntN)
				pp = append(pp, p)
			}
		}
	}

	if len(pp) == 0 {
		return nil
	}

	query := `INSERT IGNORE INTO payments(
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

	if _, err := db.Exec(query, args...); err != nil {
		return fmt.Errorf("error inserting payments: %w", err)
	}
	fmt.Println("payments created")
	return nil
}
