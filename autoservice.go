package autoservice

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Country represents countries that our services support.
type Country string

// ServiceConfig contains the configuration for individual auto services.
type ServiceConfig struct {
	Name    string
	Country Country
	Host    string
	Path    string
	Method  string
	Token   string
	Headers http.Header
}

// String returns the contents of ServiceConfig in string format.
func (cnf ServiceConfig) String() string {
	headers := ""
	for key, val := range cnf.Headers {
		headers += fmt.Sprintf("%s: %s\n", key, strings.Join(val, ","))
	}
	return fmt.Sprintf("name: %s\ncountry: %s\nhost: %s\npath: %s, method: %s\ntoken: %s\nheaders: %s\n\n", cnf.Name, cnf.Country, cnf.Host, cnf.Path, cnf.Method, cnf.Token, headers)
}

// Vehicle contains normalized vehicle data from the auto service.
type Vehicle struct {
	Brand        string
	Model        string
	RegNo        string
	VinNo        string
	FirstRegDate time.Time
}

// AutoService defines the interface that each service implementation must satisfy.
type AutoService interface {
	Configure(cnf ServiceConfig) error
	Name() string
	Supports(country Country) bool
	LookupReg(regNo string) (Vehicle, error)
	LookupVin(vinNo string) (Vehicle, error)
}

// ServiceManager keeps track of, and lets you interact with, registered auto services.
type ServiceManager map[string]AutoService

func (mngr *ServiceManager) contains(name string) bool {
	_, ok := (*mngr)[name]
	return ok
}

// AddService adds an auto service to the ServiceManager.
func (mngr *ServiceManager) AddService(service AutoService) {
	if mngr.contains(service.Name()) {
		return
	}
	(*mngr)[service.Name()] = service
}

// FindServiceByCountry returns the first service that supports the given country, or nil.
func (mngr *ServiceManager) FindServiceByCountry(country Country) AutoService {
	for _, service := range *mngr {
		if service.Supports(country) {
			return service
		}
	}
	return nil
}
