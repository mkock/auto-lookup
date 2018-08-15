package autoservice

// Service handles the shared part of various Service implementations.
type Service struct {
	Name string
	Conf ServiceConfig
}

// Configure takes a ServiceConfig (which would usually come from a YAML file).
func (service *Service) Configure(cnf ServiceConfig) error {
	service.Conf = cnf
	return nil
}

// Supports checks if the service supports license plate lookups for the given country.
func (service *Service) Supports(country Country) bool {
	return service.Conf.Country == country
}
