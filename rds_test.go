package main

import (
	"strings"
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
		"id": "au-identity-test",
		"attributes": {
			"tags.Name": "au-identity-test",
			"tags.alert": "below 10 pulse in 120",
			"tags.alert2": "below 10 swap in 60",
			"tags.alert3": "below 12 ram in 120",
			"tags.function": "operations",
			"tags.product": "identity",
			"tags.role": "au-identity-test",
			"tags.service": "database",
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
		"id": "au-identity-test",
		"attributes": {
			"tags.Name": "au-identity-test",
			"tags.alert1": "below 10 pulse in 30",
			"tags.warn1": "below 10 swap in 30",
			"tags.warn2": "below 10 cpu in 60",
			"tags.warn3": "below 10 disk in 120",
			"tags.function": "operations",
			"tags.product": "identity",
			"tags.role": "au-identity-test",
			"tags.service": "database",
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
	test1 := process(AWSDBExampleAlert, "rds")
	if len(test1) != 3 {
		t.Errorf("Should get 3 alert conditions")
	}

	test2 := process(AWSDBExampleWarn, "rds")
	if len(test2) != 4 {
		t.Errorf("Should get 4 alert conditions")
	}
}

func TestRDSParse(t *testing.T) {
	//check alert
	tag := "below 10 pulse in 30"
	duration, rule := parseCondition(strings.Fields(tag))

	if duration != 30 {
		t.Errorf("Incorrect duration %v it should be 30", duration)
	}

	if rule != "below 10 pulse" {
		t.Errorf("Incorrect alert %v it should be \"below 10 pulse\"", rule)
	}

	//Invalid tag
	tag = "below 10 pulse in as"
	duration, rule = parseCondition(strings.Fields(tag))
	if duration != 0 {
		t.Errorf("Invalid Error type")
	}

	// Invalid character count
	tag = "below 10 pulse"
	duration, rule = parseCondition(strings.Fields(tag))
	if rule != "" {
		t.Errorf("Invalid rule count")
	}

}
