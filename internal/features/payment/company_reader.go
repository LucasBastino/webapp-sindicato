package payment

import "context"

type companyReader interface {
	ListActiveIDs(ctx context.Context)([]int, error)
}

