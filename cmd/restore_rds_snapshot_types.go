package cmd

type sgRulesStruct struct {
	ruleType string `yaml:"type"`
	protocol string
	port     int
	source   string
}
