package main

import (
	"context"

	"github.com/samSRaina/hermes/services/trip-service/internal/domain"
	"github.com/samSRaina/hermes/services/trip-service/internal/infrastructure/repository"
	"github.com/samSRaina/hermes/services/trip-service/internal/service"
)

func main() {
	ctx := context.Background()
	inMemRepo := repository.NewInMemRepository()
	svc := service.NewService(inMemRepo)
	fare := &domain.RideFareModel{
		UserID: "42",
	}
	svc.CreateTrip(ctx, fare)
}
