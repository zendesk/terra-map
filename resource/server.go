package service

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type Server struct{}

type ServerMap struct {
	Server  `yaml:"-"`
	Details struct {
		ID       string `yaml:"id"`
		Alert    string `yaml:"alert,omitempty"`
		Warn     string `yaml:"warn,omitempty"`
		Duration int    `yaml:"duration"`
	} `yaml:"server"`
}

func (s Server) Process(serv *[]Service, resources []string, b []byte) {
	for _, v := range resources {
		cleanStr := strings.Replace(v, ".", "\\.", -1)

		id := fmt.Sprintf("modules.0.resources.%v.primary.attributes.tags\\.Name", cleanStr)
		itype := fmt.Sprintf("modules.0.resources.%v.primary.attributes.instance_type", cleanStr)
		resultID := gjson.Get(string(b), id)
		resultType := gjson.Get(string(b), itype)

		for _, v := range s.Conditions() {

			if strings.Contains(v.MonitorMessage, "credit") {
				if !strings.Contains(resultType.String(), "T2") ||
					!strings.Contains(resultType.String(), "T3") {
					continue
				}
			}

			server := ServerMap{}
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

func (s Server) Conditions() []Condition {

	var conditions = []Condition{
		Condition{
			MonitorType:     "alert",
			MonitorMessage:  "above 75 swap",
			MonitorDuration: 30,
		},
		Condition{
			MonitorType:     "warn",
			MonitorMessage:  "above 95 disk",
			MonitorDuration: 30,
		},
		Condition{
			MonitorType:     "creidt",
			MonitorMessage:  "below 15 credit",
			MonitorDuration: 30,
		},
		Condition{
			MonitorType:     "alert",
			MonitorMessage:  "above 95 cpu",
			MonitorDuration: 30,
		},
	}

	return conditions
}
