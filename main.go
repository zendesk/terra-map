package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	resources := getResources(dir)
	fmt.Print(string(processResources(resources)))
}

func getResources(dir string) []string {
	cmd := exec.Command("bash", "-c", "terraform show | grep -E '^[a-zA-Z]' | tr -d ':'")
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	//If repo has module then filter the name
	filteredName := strings.Replace(string(b), "module."+filepath.Base(dir)+".", "", -1)
	return strings.Split(filteredName, "\n")
}

func processResources(resources []string) (b2 []byte) {
	b, err := ioutil.ReadFile(path.Join(dir, "terraform.tfstate"))
	if err != nil {
		log.Fatal(err)
	}

	var conditions []interface{}
	for _, resource := range resources {
		if strings.Contains(resource, "aws_instance") {
			thing := Server{}
			conditions = append(conditions, thing.Process(resource, b)...)
		} else if strings.Contains(resource, "aws_sqs_queue") {
			thing := SQS{}
			conditions = append(conditions, thing.Process(resource, b)...)
		}
	}

	if len(conditions) > 0 {
		b2, err = yaml.Marshal(conditions)
		if err != nil {
			log.Fatal(err)
		}
	}
	return b2
}
