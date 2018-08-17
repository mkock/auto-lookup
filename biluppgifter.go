package autoservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// makeReq performs the request against the Biluppgifter service and the common
// part of the implementation that allows it to work with both license plate
// and VIN lookups.
func (service *BiluppgifterService) makeReq(reqURL string) (Vehicle, error) {
	// @TODO it just handles GET for now.
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return Vehicle{}, err
	}
	req.Header = service.Conf.Headers
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Vehicle{}, err
	}
	if res.StatusCode != http.StatusOK {
		return Vehicle{}, fmt.Errorf("service responded with status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Vehicle{}, err
	}
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	jsonRes := &biluppgifterData{}
	err = decoder.Decode(jsonRes)
	if err != nil {
		return Vehicle{}, err
	}
	regDate, err := time.Parse("2006-01-02", jsonRes.Data.Status.StatusData.FirstReg)
	if err != nil {
		return Vehicle{}, err
	}
	vehicle := Vehicle{
		Brand:        jsonRes.Data.Basic.Data.Make,
		Model:        jsonRes.Data.Basic.Data.Model,
		RegNo:        jsonRes.Data.Attributes.RegNo,
		VinNo:        jsonRes.Data.Attributes.VIN,
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
	return service.makeReq(reqURL.String())
}

// LookupVin looks up a vehicle based on VIN number.
func (service *BiluppgifterService) LookupVin(vinNo string) (Vehicle, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/%s/vin/%s?api_token=%s", service.Conf.Host, service.Conf.Path, vinNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	return service.makeReq(reqURL.String())
}

// Name returns the service name.
func (service *BiluppgifterService) Name() string {
	return service.Service.Name
}
