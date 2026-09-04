package creators

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func listActiveCompanyIDs(db *sqlx.DB, limit int) ([]int, error) {
	var ids []int
	err := db.Select(&ids, `
		SELECT id_company FROM companies
		WHERE deleted_at IS NULL
		ORDER BY id_company
		LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("error listing company ids for seed: %w", err)
	}
	return ids, nil
}

func formatCuilCuit(prefix int, dni string, digit int) string {
	dniNum, _ := strconv.Atoi(strings.TrimSpace(dni))
	return fmt.Sprintf("%02d-%08d-%d", prefix, dniNum, digit)
}

func generateRandomCuilCuit() string {
	middle := rand.IntN(90000000) + 10000000
	return formatCuilCuit(rand.IntN(9)+20, strconv.Itoa(middle), rand.IntN(8)+1)
}

func cuilFromDNI(dni string) string {
	return formatCuilCuit(rand.IntN(9)+20, dni, rand.IntN(8)+1)
}

func randomCompanyNumber() string {
	return strconv.Itoa(rand.IntN(900000) + 100000)
}

func randomMemberNumber() string {
	return strconv.Itoa(rand.IntN(900000000) + 100000000)
}

func isFuturePaymentPeriod(year, month int, now time.Time) bool {
	cy, cm := now.Year(), int(now.Month())
	return year > cy || (year == cy && month > cm)
}

func isPendingPaymentPeriod(year, month int, now time.Time) bool {
	cy, cm := now.Year(), int(now.Month())
	return year > cy || (year == cy && month >= cm)
}

func seedPaymentPaidAt(year, month int, now time.Time, rng func(int) int) *time.Time {
	if isPendingPaymentPeriod(year, month, now) {
		return nil
	}
	if rng(10) == 1 {
		return nil
	}
	paidAt := time.Date(year, time.Month(month), 5, 0, 0, 0, 0, time.UTC)
	return &paidAt
}

func seedPaymentAmount(year, month int, now time.Time, rng func(int) int) *float32 {
	if isPendingPaymentPeriod(year, month, now) {
		return nil
	}
	amt := float32(rng(50000) + 50000)
	return &amt
}
