package payment

import "context"

type companyReader interface {
	ListActiveIDs(ctx context.Context) ([]int, error)
	GetName(ctx context.Context, id int) (string, error)
	GetNavDetails(ctx context.Context, id int) (name, number, address, phone string, err error)
}
