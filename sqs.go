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

	id := fmt.Sprintf("modules.0.resources.%v.primary.attributes.name", cleanStr)
	alertTag := fmt.Sprintf("modules.0.resources.%v.primary.attributes.alert", cleanStr)

	var resultID gjson.Result
	if resultID = gjson.Get(string(b), id); resultID.String() == "" {
		//Module uses different path to get the data
		id = fmt.Sprintf("modules.1.resources.%v.primary.attributes.name", cleanStr)
		resultID = gjson.Get(string(b), id)
	}

	var resultAlert gjson.Result
	if resultAlert = gjson.Get(string(b), alertTag); resultAlert.String() == "" {
		//Module uses different path to get the data
		alertTag = fmt.Sprintf("modules.1.resources.%v.primary.attributes.alert", cleanStr)
		resultAlert = gjson.Get(string(b), alertTag)
	}

	if resultAlert.String() == "" {
		for _, v := range s.Conditions() {
			m := SQSCondition{}
			m.Details.Alert = v.Alert
			m.Details.Warn = v.Warn
			m.Details.ID = resultID.String()
			m.Details.Duration = v.Duration
			alerts = append(alerts, m)
		}
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
			Alert:    "above 10 visible",
			Duration: 30,
		},
		Condition{
			Alert:    "below 20 sent",
			Duration: 60,
		},
	}
}
