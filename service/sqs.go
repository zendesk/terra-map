package service

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type SQS struct{}

type SQSMap struct {
	SQS     `yaml:"-"`
	Details struct {
		ID       string `yaml:"id"`
		Alert    string `yaml:"alert,omitempty"`
		Warn     string `yaml:"warn,omitempty"`
		Duration int    `yaml:"duration"`
	} `yaml:"sqs"`
}

func (s SQS) Process(serv *[]Service, resources []string, b []byte) {
	for _, v := range resources {
		cleanStr := strings.Replace(v, ".", "\\.", -1)

		id := fmt.Sprintf("modules.0.resources.%v.primary.attributes.tags\\.Name", cleanStr)
		resultID := gjson.Get(string(b), id)

		for _, v := range s.Conditions() {
			server := SQSMap{}

			if v.MonitorType == "alert" {
				server.Details.Alert = v.MonitorMessage
			} else {
				server.Details.Warn = v.MonitorMessage
			}

			server.Details.ID = resultID.String()
			server.Details.Duration = v.MonitorDuration

			*serv = append(*serv, server)

		}

	}
}

func (s SQS) Conditions() []Condition {

	var conditions = []Condition{
		Condition{
			MonitorType:     "alert",
			MonitorMessage:  "above 75% swap",
			MonitorDuration: 30,
		},
	}

	return conditions
}
