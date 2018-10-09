package main

import (
	"fmt"
	"strings"
)

type ServerCondition struct {
	Details Condition `yaml:"server"`
}

type Server struct{}

func (s Server) Process(resource string, b []byte) (alerts []interface{}) {
	prefix := fmt.Sprintf("modules.#.resources.%v.", strings.Replace(resource, ".", "\\.", -1))
	name := queryJson(b, prefix+"primary.attributes.tags\\.Name")
	alert := queryJson(b, prefix+"primary.attributes.tags\\.alert")
	instance := queryJson(b, prefix+"primary.attributes.instance_type")

	if alert == "manual" {
		return
	}

	for _, v := range s.Conditions() {
		if strings.Contains(v.Alert, "credit") && (!strings.Contains(instance, "t2") && !strings.Contains(instance, "t3")) {
			continue
		}
		m := ServerCondition{}
		m.Details.Alert = v.Alert
		m.Details.Warn = v.Warn
		m.Details.ID = name
		m.Details.Duration = v.Duration
		alerts = append(alerts, m)
	}
	return alerts
}

func (s Server) Conditions() []Condition {
	return []Condition{
		Condition{
			Alert:    "above 75 swap",
			Duration: 30,
		},
		Condition{
			Warn:     "above 95 disk",
			Duration: 30,
		},
		Condition{
			Alert:    "below 10 credit",
			Duration: 5,
		},
		Condition{
			Warn:     "above 99 cpu",
			Duration: 120,
		},
	}
}
