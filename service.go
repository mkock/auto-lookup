package autoservice

import (
	"fmt"
	"io"
	"net/http"
)

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

// makeReq performs the request against the license plate service and returns
// the HTTP response as a byte slice.
func (service *Service) makeReq(reqURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header = service.Conf.Headers
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("service responded with status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
