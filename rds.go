package main

import (
	"strings"

	"github.com/tidwall/gjson"
)

type RDSCondition struct {
	Details Condition `yaml:"rds"`
}

type RDS struct {
	Name string
	Type string
}

func (s RDS) Process(resource string) (alerts []interface{}) {
	attr := gjson.Get(resource, "primary.attributes")
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()
	attr.ForEach(func(key, value gjson.Result) bool {
		m := RDSCondition{}
		cs := strings.Fields(value.String())
		if len(cs) == 5 {
			duration, rule := parseCondition(cs)
			if strings.Contains(key.String(), "tags.alert") {
				m.Details.Alert = rule
			} else if strings.Contains(key.String(), "tags.warn") {
				m.Details.Warn = rule
			}
			m.Details.ID = name
			m.Details.Duration = duration
			alerts = append(alerts, m)
		}
		return true
	})
	return alerts
}
