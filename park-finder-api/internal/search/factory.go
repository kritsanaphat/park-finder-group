package search

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type SearchServices struct {
	ParkingAreaStorage *storage.ParkingAreaStorage
	ReserveStorage     *storage.LogStorage
}

type ISearchServices interface {
	Search(ctx context.Context, req *models.SearchQueryRequest, keyword string) *[]models.ParkingArea
}

func NewSearchServices(
	db *mongo.Database,
) ISearchServices {
	parking_area_storage := storage.NewParkingAreaStorage(db)
	reserve_storage := storage.NewReserveStorage(db)

	return SearchServices{
		ParkingAreaStorage: parking_area_storage,
		ReserveStorage:     reserve_storage,
	}
}
