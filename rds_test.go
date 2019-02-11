package main

import (
	"fmt"
	"testing"
)

var AWSDBExampleAlert = `
{
	"type": "aws_db_instance",
	"depends_on": [
		"aws_db_subnet_group.postgres",
		"aws_security_group.postgres",
		"local.terraform"
	],
	"primary": {
		"id": "au-identity-equifax",
		"attributes": {
			"tags.Name": "au-identity-equifax",
			"tags.alert": "below 10 pulse in 30: above 2000 pulse in 30",
			"tags.function": "operations",
			"tags.product": "identity",
			"tags.role": "au-identity-equifax",
			"tags.service": "database",
			"tags.terraform": "github.com/lexerdev/ops-infra/au-identity-equifax",
		},
		"meta": {
			"e2bfb730-ecaa-11e6-8f88-34363bc7c4c0": {
				"create": 2400000000000,
				"delete": 2400000000000,
				"update": 4800000000000
			}
		},
		"tainted": false
	},
	"deposed": [],
	"provider": "provider.aws"
}
`

var AWSDBExampleWarn = `
{
	"type": "aws_db_instance",
	"depends_on": [
		"aws_db_subnet_group.postgres",
		"aws_security_group.postgres",
		"local.terraform"
	],
	"primary": {
		"id": "au-identity-equifax",
		"attributes": {
			"tags.Name": "au-identity-equifax",
			"tags.alert": "below 10 pulse in 30: above 2000 pulse in 30",
			"tags.warn": "below 10 pulse in 30",
			"tags.function": "operations",
			"tags.product": "identity",
			"tags.role": "au-identity-equifax",
			"tags.service": "database",
			"tags.terraform": "github.com/lexerdev/ops-infra/au-identity-equifax",
		},
		"meta": {
			"e2bfb730-ecaa-11e6-8f88-34363bc7c4c0": {
				"create": 2400000000000,
				"delete": 2400000000000,
				"update": 4800000000000
			}
		},
		"tainted": false
	},
	"deposed": [],
	"provider": "provider.aws"
}
`

func TestRDSProcess(t *testing.T) {
	rds := RDS{}
	test1 := rds.Process(AWSDBExampleAlert)
	if len(test1) != 2 {
		t.Errorf("Should get 2 alert conditions")
	}

	test2 := rds.Process(AWSDBExampleWarn)
	if len(test2) != 3 {
		t.Errorf("Should get 3 alert conditions")
	}
}

func TestRDSParse(t *testing.T) {
	rds := RDS{}
	rds.Type = "alert"
	rds.AlertTag = "below 10 pulse in 30: above 2000 pulse in 30"
	alerts, _ := rds.Parse(rds.AlertTag)
	if len(alerts) != 2 {
		t.Errorf("Should get 2 alert conditions")
	}

	//Invalid tag
	rds.Type = "alert"
	rds.AlertTag = "below 10 pulse in as"
	alerts, err := rds.Parse(rds.AlertTag)
	if err != nil && err.Error() != "strconv.Atoi: parsing \"as\": invalid syntax" {
		t.Errorf("Invalid Error type")
	}

	//Invalid character count
	rds.Type = "alert"
	rds.AlertTag = "below 10 pulse"
	rds.Name = "identity-uk"
	if err.Error() != fmt.Sprintf("%v alert condition needs to contain 5 characters", rds.AlertTag) {
		t.Errorf("Invalid Error type")
	}

}
