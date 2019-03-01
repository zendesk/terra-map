package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/gjson"
)

type LambdaCondition struct {
	Details Condition `yaml:"lambda"`
}

type Lambda struct {
	Name string
	Type string
}

func (s Lambda) Process(resource string) (alerts []interface{}) {
	attr := gjson.Get(resource, "primary.attributes")
	s.Name = gjson.Get(resource, "primary.attributes.function_name").String()

	attr.ForEach(func(key, value gjson.Result) bool {
		if strings.Contains(key.String(), "tags.alert") {
			s.Type = "alert"
			condition, err := s.Parse(value.String())
			if err != nil {
				log.Fatal(err)
			}
			alerts = append(alerts, condition)
		}

		if strings.Contains(key.String(), "tags.warn") {
			s.Type = "warn"
			condition, err := s.Parse(value.String())
			if err != nil {
				log.Fatal(err)
			}
			alerts = append(alerts, condition)
		}
		return true
	})

	return alerts
}

func (s Lambda) Parse(tag string) (alert interface{}, err error) {
	cs := strings.Fields(tag)
	if len(cs) == 5 {
		duration, rule, err := parseCondition(cs)
		if err != nil {
			return nil, err
		}
		m := LambdaCondition{}
		if s.Type == "alert" {
			m.Details.Alert = rule
		} else if s.Type == "warn" {
			m.Details.Warn = rule
		}
		m.Details.ID = s.Name
		m.Details.Duration = duration
		return m, nil
	} else {
		return nil, fmt.Errorf("%v alert condition needs to contain 5 words", tag)
	}
}
