package autoservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// BiluppgifterService integrates with a Danish license plate lookup service.
type BiluppgifterService struct {
	Service
}

// biluppgifterData represents the HTTP response from biluppgifter.
type biluppgifterData struct {
	Data biluppgifterResponse `json:"data"`
}

// biluppgifterResponse represents the main HTTP response (inside "data") from biluppgifter.
type biluppgifterResponse struct {
	Attributes attributeResponse `json:"attributes"`
	Basic      basicResponse     `json:"basic"`
	Status     statusResponse    `json:"status"`
}

type basicResponse struct {
	Data basicDataResponse `json:"data"`
}

type basicDataResponse struct {
	Make  string `json:"make"`
	Model string `json:"model"`
}

type statusResponse struct {
	StatusData statusDataResponse `json:"data"`
}

type statusDataResponse struct {
	FirstReg string `json:"first_registered"`
}

// attributeResponse contains the attributes from the response.
type attributeResponse struct {
	RegNo string `json:"regno"`
	VIN   string `json:"vin"`
}

// decodeBiluppgifterBody decodes the HTTP response into a Vehicle.
func decodeBiluppgifterBody(body []byte) (Vehicle, error) {
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	data := &biluppgifterData{}
	err := decoder.Decode(data)
	if err != nil {
		return Vehicle{}, err
	}
	regDate, err := time.Parse("2006-01-02", data.Data.Status.StatusData.FirstReg)
	if err != nil {
		return Vehicle{}, err
	}
	vehicle := Vehicle{
		Brand:        data.Data.Basic.Data.Make,
		Model:        data.Data.Basic.Data.Model,
		RegNo:        data.Data.Attributes.RegNo,
		VinNo:        data.Data.Attributes.VIN,
		FirstRegDate: regDate,
	}
	return vehicle, nil
}

// LookupReg looks up a vehicle based on registration number.
func (service *BiluppgifterService) LookupReg(regNo string) (Vehicle, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s/regno/%s?api_token=%s", service.Conf.Host, service.Conf.Path, regNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	body, err := service.makeReq(reqURL.String())
	if err != nil {
		return Vehicle{}, err
	}
	return decodeBiluppgifterBody(body)
}

// LookupVin looks up a vehicle based on VIN number.
func (service *BiluppgifterService) LookupVin(vinNo string) (Vehicle, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s/vin/%s?api_token=%s", service.Conf.Host, service.Conf.Path, vinNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	body, err := service.makeReq(reqURL.String())
	if err != nil {
		return Vehicle{}, err
	}
	return decodeBiluppgifterBody(body)
}

// Name returns the service name.
func (service *BiluppgifterService) Name() string {
	return service.Service.Name
}
