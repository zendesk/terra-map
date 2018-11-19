package main

import "github.com/tidwall/gjson"

type SQSCondition struct {
	Details Condition `yaml:"sqs"`
}

type SQS struct{}

func (s SQS) Process(resource string) (alerts []interface{}) {
	name := gjson.Get(resource, "primary.attributes.name").String()
	alert := gjson.Get(resource, "primary.attributes.tags\\.alert").String()

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
