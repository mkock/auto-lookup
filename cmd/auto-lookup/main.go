package main

import (
	"fmt"
	"net/http"

	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/mkock/auto-lookup"
)

func readConfigFrom(fname string) ([]autoservice.ServiceConfig, error) {
	conf, err := yaml.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	nodes, err := yaml.Child(conf.Root, "services")
	services, servicesOk := nodes.(yaml.Map)
	if !servicesOk {
		return nil, fmt.Errorf("parse YAML file: section 'services' is not a list")
	}
	configs := make([]autoservice.ServiceConfig, 0, len(services))
	var (
		opts    yaml.Map
		headers yaml.Map
		ok      bool
	)
	for name, val := range services {
		fmt.Printf("Identified service %q\n", name)
		if opts, ok = val.(yaml.Map); !ok {
			return nil, fmt.Errorf("parse YAML file: service %q does not contain a list of options", name)
		}
		node, err := yaml.Child(opts, ".headers")
		if err != nil {
			return nil, fmt.Errorf("parse YAML file: service %q contains invalid headers", name)
		}
		headers = node.(yaml.Map)
		httpHeaders := make(http.Header, len(services))
		for k, v := range headers {
			httpHeaders.Add(k, v.(yaml.Scalar).String())
		}
		country := autoservice.Country(opts.Key("country").(yaml.Scalar).String())
		config := autoservice.ServiceConfig{
			Name:    name,
			Country: country,
			Host:    opts.Key("host").(yaml.Scalar).String(),
			Path:    opts.Key("path").(yaml.Scalar).String(),
			Method:  opts.Key("method").(yaml.Scalar).String(),
			Token:   opts.Key("token").(yaml.Scalar).String(),
			Headers: httpHeaders,
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func main() {
	configs, err := readConfigFrom("conf.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(configs); i++ {
		fmt.Printf("%v", configs[i])
	}
}
