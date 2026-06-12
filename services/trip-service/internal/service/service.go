package service

import (
	"context"

	"github.com/samSRaina/hermes/services/trip-service/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       bson.NewObjectID(),
		UserID:   "",
		Status:   "",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)
}
