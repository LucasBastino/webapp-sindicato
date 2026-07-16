package idempotency

type IdempotencyHandler struct {
	service *IdempotencyService
}

func NewIdempotencyHandler(service *IdempotencyService) *IdempotencyHandler {
	return &IdempotencyHandler{service: service}
}

