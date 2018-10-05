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

	b, err := processServices(resourceMap)
	if err != nil {
		log.Fatal(err)
	}

	err = appendDateToMap(b)
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

func processServices(resourceMap map[string][]string) ([]byte, error) {

	b, err := ioutil.ReadFile(path.Join(dir, "terraform.tfstate"))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return b2, nil
}

func appendDateToMap(b []byte) error {

	fileName := path.Join(dir, "map.yml")

	if _, err := os.Stat(path.Join(fileName)); err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", "sed -i '/# automatically created alerts below this/q' "+fileName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(string(b)); err != nil {
		return err
	}

	return nil
}
