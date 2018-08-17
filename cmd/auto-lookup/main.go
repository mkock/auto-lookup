package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mkock/auto-lookup"
)

func initServices() []autoservice.AutoService {
	nrpladeService := autoservice.NrpladeService{Service: autoservice.Service{Name: "nrplade"}}
	biluppgifterService := autoservice.BiluppgifterService{Service: autoservice.Service{Name: "biluppgifter"}}
	services := []autoservice.AutoService{
		&nrpladeService,
		&biluppgifterService,
	}
	return services
}

func main() {
	// Get regno/vin no from CLI.
	if len(os.Args) < 3 {
		fmt.Println("Please follow program name by a country and registration or VIN number to lookup, eg. 'auto-lookup.exe dk BX71743'")
		fmt.Println("Or, alternatively, write 'test' as registration number to lookup a hard-coded (Danish) registration number.")
		return
	}

	configs, err := autoservice.ReadConfigFrom("conf.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	mngr := autoservice.ServiceManager{}

	// Configure each service dynamically.
	services := initServices()
	for _, service := range services {
		for i := 0; i < len(configs); i++ {
			if configs[i].Name == service.Name() {
				service.Configure(configs[i])
				mngr.AddService(service)
			}
		}
	}
	fmt.Printf("Successfully registered %d services.\n", len(mngr))

	// Find and call the appropriate service.
	country := os.Args[1]
	if country == "" {
		fmt.Println("Missing country identifier: dk or se.")
		return
	}
	service := mngr.FindServiceByCountry(autoservice.Country(country))
	if service == nil {
		fmt.Printf("No service available for country: %s.\n", country)
		os.Exit(1)
	}
	regNo := os.Args[2]
	if strings.EqualFold(regNo, "test") {
		regNo = "BX71743" // Test reg.
	}
	fmt.Printf("Looking for %s... ", regNo)
	var vehicle autoservice.Vehicle
	if autoservice.IsVIN(regNo) {
		vehicle, err = service.LookupVin(regNo)
	} else {
		vehicle, err = service.LookupReg(regNo)
	}
	if err != nil {
		fmt.Printf("lookup by regno: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("found!")
	fmt.Printf("%s\n", vehicle)
}
