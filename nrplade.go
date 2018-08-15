package autoservice

import (
	"time"
)

// NrpladeService integrates with a Danish license plate lookup service.
type NrpladeService struct {
	Service
}

// LookupReg looks up a vehicle based on registration number.
func (service *NrpladeService) LookupReg(regNo string) (Vehicle, error) {
	vehicle := Vehicle{
		Brand:        "Ford",
		Model:        "GT",
		RegNo:        "123",
		VinNo:        "123",
		FirstRegDate: time.Now(),
	}
	return vehicle, nil
}

// LookupVin looks up a vehicle based on VIN number.
func (service *NrpladeService) LookupVin(vinNo string) (Vehicle, error) {
	vehicle := Vehicle{
		Brand:        "Ford",
		Model:        "GT",
		RegNo:        "123",
		VinNo:        "123",
		FirstRegDate: time.Now(),
	}
	return vehicle, nil
}

// Name returns the service name.
func (service *NrpladeService) Name() string {
	return service.Service.Name
}
