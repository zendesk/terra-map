package main

import (
	"io/ioutil"
	"path"
	"sort"

	"github.com/tidwall/gjson"
	yaml "gopkg.in/yaml.v2"
)

type AppCondition struct {
	Details Condition `yaml:"pulse"`
}

type App struct{}

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

func (s App) Process(resource string) (alerts []interface{}) {
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()

	services := getServices()
	sort.Strings(services)

	for _, service := range services {
		for _, v := range s.Conditions() {
			m := AppCondition{}
			m.Details.Alert = v.Alert
			m.Details.Warn = v.Warn
			m.Details.ID = name + "/" + service
			m.Details.Duration = v.Duration
			alerts = append(alerts, m)
		}
	}

	return alerts
}

func (s App) Conditions() []Condition {
	return []Condition{
		Condition{
			Warn:     "below 5 pulse",
			Duration: 30,
		},
	}
}
