package autoservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// NrpladeService integrates with a Danish license plate lookup service.
type NrpladeService struct {
	Service
}

// nrpladeData represents the HTTP response from nrpla.de.
type nrPladeData struct {
	Data nrpladeResponse `json:"data"`
}

// nrpladeResponse represents the main HTTP response (inside "data") from nrpla.de.
type nrpladeResponse struct {
	Registration string `json:"registration"`
	FirstRegDate string `json:"first_registration_date"`
	Vin          string `json:"vin"`
	Type         string `json:"type"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Version      string `json:"version"`
	FuelType     string `json:"fuel_type"`
	RegStatus    string `json:"registration_status"`
}

// decodeNrpladeBody does the decoding of the HTTP response into Vehicle.
func decodeNrpladeBody(body []byte) (Vehicle, error) {
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	data := &nrPladeData{}
	err := decoder.Decode(data)
	if err != nil {
		return Vehicle{}, err
	}
	regDate, err := time.Parse("2006-01-02", data.Data.FirstRegDate)
	if err != nil {
		return Vehicle{}, err
	}
	vehicle := Vehicle{
		Brand:        data.Data.Brand,
		Model:        data.Data.Model,
		RegNo:        data.Data.Registration,
		VinNo:        data.Data.Vin,
		FirstRegDate: regDate,
	}
	return vehicle, nil
}

// LookupReg looks up a vehicle based on registration number.
func (service *NrpladeService) LookupReg(regNo string) (Vehicle, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s/%s?api_token=%s", service.Conf.Host, service.Conf.Path, regNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	body, err := service.makeReq(reqURL.String())
	if err != nil {
		return Vehicle{}, err
	}
	return decodeNrpladeBody(body)
}

// LookupVin looks up a vehicle based on VIN number.
func (service *NrpladeService) LookupVin(vinNo string) (Vehicle, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s/vin/%s?api_token=%s", service.Conf.Host, service.Conf.Path, vinNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	body, err := service.makeReq(reqURL.String())
	if err != nil {
		return Vehicle{}, err
	}
	return decodeNrpladeBody(body)
}

// Name returns the service name.
func (service *NrpladeService) Name() string {
	return service.Service.Name
}
