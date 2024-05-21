package service

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/repository"
)

type LocationService struct {
	r *repository.LocationRepository
}

func (l *LocationService) Insert(location model.InsertLocationForm) error {
	return l.r.Insert(location)
}

func (l *LocationService) Nearest(long, lat float64, maxDistance int, limit int) ([]model.Location, error) {
	return l.r.Nearest(long, lat, maxDistance, limit)
}

func NewLocationService(r *repository.LocationRepository) *LocationService {
	return &LocationService{r}
}
