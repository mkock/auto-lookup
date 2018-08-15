package main

import (
	"fmt"
	"os"

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
	regNo := "RL90123"
	vehicle, err := service.LookupReg(regNo)
	if err != nil {
		fmt.Printf("lookup by regno: %s", err.Error())
		os.Exit(0)
	}
	fmt.Println("found!")
	fmt.Println(vehicle)
}
