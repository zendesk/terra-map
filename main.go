package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

var dir string

type Resource interface {
	Process(resource []string, b []byte) []interface{}
	Conditions() []Condition
}

type Condition struct {
	ID       string `yaml:"id"`
	Alert    string `yaml:"alert,omitempty"`
	Warn     string `yaml:"warn,omitempty"`
	Duration int    `yaml:"duration"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s DIR", os.Args[0])
	}
	dir := os.Args[1]
	if _, err := os.Stat(path.Join(dir, "terraform.tfstate")); err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadFile(path.Join(dir, "terraform.tfstate"))
	if err != nil {
		log.Fatal(err)
	}

	resources := getResources(string(b))
	fmt.Print(string(processResources(string(b), resources)))
}

func queryJson(state string, search string) (val string) {
	if result := gjson.Get(state, search).Array(); len(result) > 0 {
		val = result[0].String()
	}
	return val
}

func getResources(state string) (resources []string) {
	if result := gjson.Get(state, "modules.#.resources").Array(); len(result) > 0 {
		for _, v := range result {
			for k, _ := range v.Map() {
				resources = append(resources, k)
			}
		}
	}
	sort.Strings(resources)
	return resources
}

func processResources(state string, resources []string) (b2 []byte) {
	var conditions []interface{}
	for _, resource := range resources {
		if strings.Contains(resource, "aws_instance") {
			thing := Server{}
			conditions = append(conditions, thing.Process(state, resource)...)
		} else if strings.Contains(resource, "aws_sqs_queue") {
			thing := SQS{}
			conditions = append(conditions, thing.Process(state, resource)...)
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
