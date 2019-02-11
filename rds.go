package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
)

type RDSCondition struct {
	Details Condition `yaml:"rds"`
}

type RDS struct {
	Name     string
	Type     string
	AlertTag string
	WarnTag  string
}

func (s RDS) Process(resource string) (alerts []interface{}) {
	s.Name = gjson.Get(resource, "primary.attributes.tags\\.Name").String()
	s.AlertTag = gjson.Get(resource, "primary.attributes.tags\\.alert").String()
	s.WarnTag = gjson.Get(resource, "primary.attributes.tags\\.warn").String()

	if s.AlertTag != "" {
		s.Type = "alert"
		conditions, err := s.Parse(s.AlertTag)
		if err != nil {
			log.Fatal(err)
		}
		alerts = append(alerts, conditions...)
	}

	if s.WarnTag != "" {
		s.Type = "warn"
		conditions, err := s.Parse(s.WarnTag)
		if err != nil {
			log.Fatal(err)
		}
		alerts = append(alerts, conditions...)
	}
	return alerts
}

func (s RDS) Parse(tag string) (alerts []interface{}, err error) {
	conditions := strings.Split(strings.TrimSpace(tag), ":")
	sort.Strings(conditions)
	for _, c := range conditions {
		cs := strings.Fields(c)
		if len(cs) == 5 {
			duration, rule, err := parseCondition(cs)
			if err != nil {
				return nil, err
			}
			m := RDSCondition{}
			if s.Type == "alert" {
				m.Details.Alert = rule
			} else if s.Type == "warn" {
				m.Details.Warn = rule
			}
			m.Details.ID = s.Name
			m.Details.Duration = duration
			alerts = append(alerts, m)
		} else {
			return nil, fmt.Errorf("%v alert condition needs to contain 5 characters", c)
		}
	}
	return alerts, nil
}
