package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/gjson"
)

type ESCondition struct {
	Details Condition `yaml:"es"`
}

type ES struct {
	Name string
	Type string
}

func (s ES) Process(resource string) (alerts []interface{}) {
	attr := gjson.Get(resource, "primary.attributes")
	s.Name = gjson.Get(resource, "primary.id").String()

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

func (s ES) Parse(tag string) (alert interface{}, err error) {
	cs := strings.Fields(tag)
	if len(cs) == 5 {
		duration, rule := parseCondition(cs)
		m := ESCondition{}
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
