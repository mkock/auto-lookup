package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mkock/auto-lookup"
)

func initServices() []autoservice.AutoService {
	nrpladeService := autoservice.NrpladeService{Service: autoservice.Service{Name: "nrplade"}}
	services := []autoservice.AutoService{
		&nrpladeService,
	}
	return services
}

func main() {
	// Get regno/vin no from CLI.
	if len(os.Args) < 2 {
		fmt.Println("Please follow program name by a registration or VIN number to lookup, eg. 'auto-lookup.exe BX71743'")
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
			fmt.Printf("%v", configs[i])
			if configs[i].Name == service.Name() {
				service.Configure(configs[i])
				mngr.AddService(service)
			}
		}
	}
	fmt.Printf("Successfully registered %d services.\n", len(mngr))

	// Find and call a Danish service.
	service := mngr.FindServiceByCountry(autoservice.Country("dk"))
	if service == nil {
		fmt.Println("No service available for country dk.")
		os.Exit(1)
	}
	regNo := os.Args[1]
	if strings.EqualFold(regNo, "test") {
		regNo = "BX71743" // Test reg.
	}
	fmt.Printf("Looking for %s...\n", regNo)
	vehicle, err := service.LookupReg(regNo)
	if err != nil {
		fmt.Printf("lookup by regno: %s", err.Error())
		os.Exit(0)
	}
	fmt.Println("found!")
	fmt.Printf("%s\n", vehicle)
}
