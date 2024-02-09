package cmd

type sgRulesStruct struct {
	ruleType string `yaml:"type"`
	protocol string `yaml:"protocol"`
	port     int    `yaml:"port"`
	source   string `yaml:"source"`
}

type sgStruct struct {
	snapshot          string
	database          string
	region            string
	dbsubnetgroupname string
	subnets           []string
	securitygroups    []sgRulesStruct
}
