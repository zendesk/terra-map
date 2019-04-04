package main

import (
	"io/ioutil"
	"path"
	"sort"

	"github.com/tidwall/gjson"
	yaml "gopkg.in/yaml.v2"
)

// Pulse to model docker services
type Pulse struct{}

// DockerCompose to model a docker-compose.yml file
type DockerCompose struct {
	Version  string `yaml:"version"`
	Services map[string]struct {
		ContainerName string            `yaml:"container_name"`
		Image         string            `yaml:"image"`
		Labels        map[string]string `yaml:"labels"`
	} `yaml:"services"`
}

func getServices() []string {
	var services []string

	dc, err := ioutil.ReadFile(path.Join(dir, "docker-compose.yml"))
	if err != nil {
		return services
	}

	structure := DockerCompose{}
	err = yaml.Unmarshal(dc, &structure)
	if err != nil {
		return services
	}

	for key := range structure.Services {
		if structure.Services[key].Labels["alert"] == "manual" {
			continue
		}
		services = append(services, structure.Services[key].ContainerName)
	}

	return services
}

func (s Pulse) process(resource string) (alerts []condition) {
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()

	services := getServices()
	sort.Strings(services)

	alerts = []condition{}
	for _, service := range services {
		con := condition{"pulse": {ID: name + "/" + service, Alert: "", Warn: "below 5 pulse", Duration: 120}}
		alerts = append(alerts, con)
	}

	return alerts
}
