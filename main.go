package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/shoukoo/terra-map/service"
	"gopkg.in/yaml.v2"
)

var dir string
var validResource = []string{"aws_instance", "aws_s3_bucket", "aws_sqs_queue"}

func init() {
	flag.StringVar(&dir, "d", "", "specify a dir path to generate a map.yml")
	flag.Parse()
}

func main() {

	if dir == "" {
		log.Fatalf("Provide a dir path e.g. -d=/home/andy/ops/ ")
	}

	//check if terraform.tfstate exists in this folder
	if _, err := os.Stat(path.Join(dir, "terraform.tfstate")); err != nil {
		log.Fatal(err)
	}

	//cd to that dir
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	resourceMap, err := getListOfResources()
	if err != nil {
		log.Fatal(err)
	}

	err = processServices(resourceMap)
	if err != nil {
		log.Fatal(err)
	}
}

func getListOfResources() (map[string][]string, error) {

	resourceMap := make(map[string][]string)

	cmd := exec.Command("bash", "-c", "terraform show | grep -E '^[a-zA-Z]' | tr -d ':'")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	splitOutput := strings.Split(string(b), "\n")
	for _, v := range splitOutput {
		for _, r := range validResource {
			if strings.Contains(v, r) {
				resourceMap[r] = append(resourceMap[r], v)
			}
		}
	}

	return resourceMap, nil
}

func processServices(resourceMap map[string][]string) error {

	b, err := ioutil.ReadFile(path.Join(dir, "terraform.tfstate"))
	if err != nil {
		return err
	}

	var servers = []service.Service{}

	if len(resourceMap["aws_instance"]) > 0 {
		server := service.Server{}
		server.Process(&servers, resourceMap["aws_instance"], b)
	}

	if len(resourceMap["aws_sqs_queue"]) > 0 {
		sqs := service.SQS{}
		sqs.Process(&servers, resourceMap["aws_sqs_queue"], b)
	}

	if len(resourceMap["aws_s3_bucket"]) > 0 {
		sqs := service.SQS{}
		sqs.Process(&servers, resourceMap["aws_s3_bucket"], b)
	}

	b2, err := yaml.Marshal(servers)
	if err != nil {
		return err
	}

	//We can changhe the file name later...
	err = ioutil.WriteFile("terra_map.yml", b2, 0777)
	if err != nil {
		return err
	}

	return nil
}
