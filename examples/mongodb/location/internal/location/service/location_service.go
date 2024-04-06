package service

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/repository"
)

type LocationService interface {
	Insert(location model.InsertLocationForm) error
	Nearest(long, lat float64, maxDistance int, limit int) ([]model.Location, error)
}

type locationService struct {
	r repository.LocationRepository
}

func (l locationService) Insert(location model.InsertLocationForm) error {
	return l.r.Insert(location)
}

func (l locationService) Nearest(long, lat float64, maxDistance int, limit int) ([]model.Location, error) {
	return l.r.Nearest(long, lat, maxDistance, limit)
}

func NewLocationService(r repository.LocationRepository) LocationService {
	return &locationService{r}
}
