package main

import (
	"io/ioutil"
	"path"
	"strings"

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

func (s Pulse) process(resource string) (alerts []condition) {

	// get name of server
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()

	// examine docker-compose.yml for alert conditions defined in docker labels
	dc, err := ioutil.ReadFile(path.Join(dir, "docker-compose.yml"))
	if err != nil {
		return
	}

	structure := DockerCompose{}
	err = yaml.Unmarshal(dc, &structure)
	if err != nil {
		return
	}

	alerts = []condition{}
	for key := range structure.Services {
		container := structure.Services[key].ContainerName
		for k, v := range structure.Services[key].Labels {
			if duration, rule := parseCondition(v); duration != 0 && rule != "" {
				if strings.Contains(k, "alert") {
					alerts = append(alerts, condition{"pulse": {ID: name + "/" + container, Alert: rule, Duration: duration}})
				} else if strings.Contains(k, "warn") {
					alerts = append(alerts, condition{"pulse": {ID: name + "/" + container, Warn: rule, Duration: duration}})
				}
			}
		}
	}

	return alerts
}
