package main

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type SQSCondition struct {
	Details Condition `yaml:"sqs"`
}

type SQS struct{}

func (s SQS) Process(resource string, b []byte) (alerts []interface{}) {
	cleanStr := strings.Replace(resource, ".", "\\.", -1)
	id := fmt.Sprintf("modules.0.resources.%v.primary.attributes.tags\\.Name", cleanStr)
	resultID := gjson.Get(string(b), id)
	for _, v := range s.Conditions() {
		m := SQSCondition{}
		m.Details.Alert = v.Alert
		m.Details.Warn = v.Warn
		m.Details.ID = resultID.String()
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
		Condition{
			Alert:    "below 20 sent",
			Duration: 60,
		},
	}
}