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

	var resultID gjson.Result
	if resultID = gjson.Get(string(b), id); resultID.String() == "" {
		//Module uses different path to get the data
		id = fmt.Sprintf("modules.1.resources.%v.primary.attributes.name", cleanStr)
		resultID = gjson.Get(string(b), id)
	}

	for _, v := range s.Conditions() {

		if v.Pattern == "" || (v.Pattern != "" && processPattern(resultID.String(), v.Pattern)) {
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
			Pattern:  "not dead",
			Duration: 60,
		},
		Condition{
			Alert:    "above 10 visible",
			Pattern:  "equal dead",
			Duration: 30,
		},
		Condition{
			Alert:    "below 20 sent",
			Duration: 60,
		},
	}
}
