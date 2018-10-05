package main

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type ServerCondition struct {
	Details Condition `yaml:"server"`
}

type Server struct{}

func (s Server) Process(resource string, b []byte) (alerts []interface{}) {
	cleanStr := strings.Replace(resource, ".", "\\.", -1)

	id := fmt.Sprintf("modules.0.resources.%v.primary.attributes.tags\\.Name", cleanStr)
	itype := fmt.Sprintf("modules.0.resources.%v.primary.attributes.instance_type", cleanStr)
	resultID := gjson.Get(string(b), id)
	resultType := gjson.Get(string(b), itype)

	for _, v := range s.Conditions() {
		if strings.Contains(v.Alert, "credit") {
			if !strings.Contains(resultType.String(), "t2") &&
				!strings.Contains(resultType.String(), "t3") {
				continue
			}
		}

		m := ServerCondition{}
		m.Details.Alert = v.Alert
		m.Details.Warn = v.Warn
		m.Details.ID = resultID.String()
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
