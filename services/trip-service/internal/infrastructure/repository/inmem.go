package repository

import (
	"context"
	"hermes/services/trip-service/internal/domain"
)

type inMemRepository struct {
	trips    map[string]*domain.TripModel
	rideFare map[string]*domain.RideFareModel
}

func NewInMemRepository() *inMemRepository {
	return &inMemRepository{
		trips:    make(map[string]*domain.TripModel),
		rideFare: make(map[string]*domain.RideFareModel),
	}
}

func (r *inMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
