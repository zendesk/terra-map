package main

import (
	"fmt"
	"strings"
)

type AppCondition struct {
	Details Condition `yaml:"app"`
}

type App struct{}

func (s App) Process(state string, resource string, options ...interface{}) (alerts []interface{}) {
	prefix := fmt.Sprintf("modules.#.resources.%v.", strings.Replace(resource, ".", "\\.", -1))
	name := queryJson(state, prefix+"primary.attributes.tags\\.Name")

	var services []string
	if len(options) > 0 {
		services = options[0].([]string)
	}

	fmt.Println(services)

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
			Alert:    "below 5 pulse",
			Duration: 30,
		},
	}
}
