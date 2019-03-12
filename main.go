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

// Condition alert conditions
type Condition struct {
	ID       string `yaml:"id"`
	Alert    string `yaml:"alert,omitempty"`
	Warn     string `yaml:"warn,omitempty"`
	Duration int    `yaml:"duration"`
}

func main() {

	if len(os.Args) != 2 {
		log.Println("Version: v2.5.2")
		log.Fatalf("Usage: %s DIR", os.Args[0])
	}

	dir = os.Args[1]
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

func processResources(resources []string) (b2 []byte) {
	var conditions []interface{}
	for _, resource := range resources {
		if gjson.Get(resource, "type").String() == "aws_instance" {

			server := Server{}
			conditions = append(conditions, server.Process(resource)...)

			app := App{}
			conditions = append(conditions, app.Process(resource)...)

		} else if gjson.Get(resource, "type").String() == "aws_sqs_queue" {
			sqs := SQS{}
			conditions = append(conditions, sqs.Process(resource)...)
		} else if gjson.Get(resource, "type").String() == "aws_lambda_function" {
			lambda := Lambda{}
			conditions = append(conditions, lambda.Process(resource)...)
		} else if gjson.Get(resource, "type").String() == "aws_db_instance" {
			rds := RDS{}
			conditions = append(conditions, rds.Process(resource)...)
		}
	}

	if len(conditions) > 0 {
		b2, err := yaml.Marshal(conditions)
		if err != nil {
			log.Fatal(err)
		} else {
			return b2
		}
	}
	return b2
}

func parseCondition(conditon []string) (duration int, rule string, err error) {
	duration, err = strconv.Atoi(strings.Join(conditon[len(conditon)-1:], " "))
	if err != nil {
		return 0, "", err
	}
	rule = strings.Join(conditon[:len(conditon)-2], " ")
	return duration, rule, nil
}
