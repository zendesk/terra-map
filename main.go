package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

var dir string

type condition map[string]struct {
	ID       string `yaml:"id"`
	Alert    string `yaml:"alert,omitempty"`
	Warn     string `yaml:"warn,omitempty"`
	Duration int    `yaml:"duration"`
}

func main() {
	log.SetPrefix("terra-map v2.5.3 ")
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "."
	}
	if _, err := os.Stat(path.Join(dir, "terraform.tfstate")); err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadFile(path.Join(dir, "terraform.tfstate"))
	if err != nil {
		log.Fatal(err)
	}

	resources := getResources(string(b))
	fmt.Print(string(processResources(resources)))
}

func getResources(state string) (resources []string) {
	if result := gjson.Get(state, "modules.#.resources").Array(); len(result) > 0 {
		for _, v := range result {

			var keys []string
			for k := range v.Map() {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				// ignore all data resource
				if !strings.HasPrefix(k, "data.") {
					resources = append(resources, v.Map()[k].Raw)
				}
			}
		}
	}
	return resources
}

func processResources(resources []string) (b []byte) {
	var conditions []condition
	for _, resource := range resources {
		if gjson.Get(resource, "type").String() == "aws_instance" {
			conditions = append(conditions, process(resource, "server")...)

			// special case where we need to parse a docker-compose.yml file
			pulse := Pulse{}
			conditions = append(conditions, pulse.process(resource)...)

		} else if gjson.Get(resource, "type").String() == "aws_sqs_queue" {
			conditions = append(conditions, process(resource, "sqs")...)

		} else if gjson.Get(resource, "type").String() == "aws_lambda_function" {
			conditions = append(conditions, process(resource, "lambda")...)

		} else if gjson.Get(resource, "type").String() == "aws_db_instance" {
			conditions = append(conditions, process(resource, "rds")...)

		} else if gjson.Get(resource, "type").String() == "aws_ssm_parameter" {
			conditions = append(conditions, process(resource, "es")...)
		}
	}

	if len(conditions) > 0 {
		b, err := yaml.Marshal(conditions)
		if err != nil {
			log.Fatal(err)
		} else {
			return b
		}
	}
	return
}

func parseCondition(condition []string) (duration int, rule string) {
	duration, err := strconv.Atoi(strings.Join(condition[len(condition)-1:], " "))
	if err == nil {
		rule = strings.Join(condition[:len(condition)-2], " ")
		return duration, rule
	}
	return
}

func process(resource string, thing string) (alerts []condition) {
	attr := gjson.Get(resource, "primary.attributes")
	name := gjson.Get(resource, "primary.attributes.tags\\.Name").String()
	if name == "" {
		name = gjson.Get(resource, "primary.attributes.function_name").String()
	}
	if name == "" {
		name = gjson.Get(resource, "primary.id").String()
	}
	if name == "" {
		name = gjson.Get(resource, "primary.attributes.id").String()
	}

	alerts = []condition{}
	attr.ForEach(func(key, value gjson.Result) bool {
		cs := strings.Fields(value.String())
		if len(cs) == 5 && strings.Contains(key.String(), "tags.alert") {
			duration, rule := parseCondition(cs)
			con := condition{thing: {ID: name, Alert: rule, Duration: duration}}
			alerts = append(alerts, con)

		} else if len(cs) == 5 && strings.Contains(key.String(), "tags.warn") {
			duration, rule := parseCondition(cs)
			con := condition{thing: {ID: name, Warn: rule, Duration: duration}}
			alerts = append(alerts, con)
		}
		return true
	})
	return alerts
}
