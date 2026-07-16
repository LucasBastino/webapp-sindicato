package idempotency

import (
	"context"
	"errors"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
)

type IdempotencyService struct {
	repo *IdempotencyRepository
}

func NewIdempotencyService(repo *IdempotencyRepository, ttl time.Duration) *IdempotencyService {
	return &IdempotencyService{
		repo: repo,
	}
}

func (s *IdempotencyService) CheckOrCreate(ctx context.Context, key string, requestHash string) (*Record, error) {

	record, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if record != nil {

		if record.RequestHash != requestHash {
			return nil, apperrors.NewBusinessError(errors.New("idempotency key reused with different request"), "")
		}

		return record, nil
	}

	err = s.repo.Create(ctx, key, requestHash, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return nil, nil
}

func (s *IdempotencyService) UpdateResource(ctx context.Context, key string, resourceType string, resourceID int) error {
	err := s.repo.UpdateResource(ctx, key, resourceType, resourceID)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *IdempotencyService) CleanupExpiredKeys(ctx context.Context) error {
	if err := s.repo.DeleteExpired(ctx); err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}