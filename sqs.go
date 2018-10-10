package main

import (
	"fmt"
	"strings"
)

type SQSCondition struct {
	Details Condition `yaml:"sqs"`
}

type SQS struct{}

func (s SQS) Process(state string, resource string) (alerts []interface{}) {
	prefix := fmt.Sprintf("modules.#.resources.%v.", strings.Replace(resource, ".", "\\.", -1))
	name := queryJson(state, prefix+"primary.attributes.name")
	alert := queryJson(state, prefix+"primary.attributes.tags\\.alert")

	if alert == "manual" {
		return
	}

	for _, v := range s.Conditions() {
		m := SQSCondition{}
		m.Details.Alert = v.Alert
		m.Details.Warn = v.Warn
		m.Details.ID = name
		m.Details.Duration = v.Duration
		alerts = append(alerts, m)
	}
	return alerts
}

func (s SQS) Conditions() []Condition {
	return []Condition{
		Condition{
			Alert:    "above 5000 visible",
			Duration: 60,
		},
	}
}
