package installment

import "context"

type paymentPlanStatusRefresher interface {
	RefreshStatus(ctx context.Context, id int) error
	GetStatus(ctx context.Context, id int) (string, error)
}
