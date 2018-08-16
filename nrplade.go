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

// LookupReg looks up a vehicle based on registration number.
func (service *NrpladeService) LookupReg(regNo string) (Vehicle, error) {
	// @TODO it just handles GET for now.
	u, err := url.Parse(fmt.Sprintf("%s/%s/%s?api_token=%s", service.Conf.Host, service.Conf.Path, regNo, service.Conf.Token))
	if err != nil {
		return Vehicle{}, err
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Vehicle{}, err
	}
	req.Header = service.Conf.Headers
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Vehicle{}, err
	}
	fmt.Println(req.URL)
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
	jsonRes := &nrPladeData{}
	err = decoder.Decode(jsonRes)
	if err != nil {
		return Vehicle{}, err
	}
	regDate, err := time.Parse("2006-01-02", jsonRes.Data.FirstRegDate)
	if err != nil {
		return Vehicle{}, err
	}
	vehicle := Vehicle{
		Brand:        jsonRes.Data.Brand,
		Model:        jsonRes.Data.Model,
		RegNo:        jsonRes.Data.Registration,
		VinNo:        jsonRes.Data.Vin,
		FirstRegDate: regDate,
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
