package main

import (
	"strings"

	"github.com/tidwall/gjson"
)

type ServerCondition struct {
	Details Condition `yaml:"server"`
}

type Server struct{}

func (s Server) Process(resource string) (alerts []interface{}) {
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()
	alert := gjson.Get(resource, "primary.attributes.tags\\.alert").String()
	instance := gjson.Get(resource, "primary.attributes.instance_type").String()

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
