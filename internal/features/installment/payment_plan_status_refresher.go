package installment

import "context"

type paymentPlanStatusRefresher interface {
	RefreshStatus(ctx context.Context, id int) error
}