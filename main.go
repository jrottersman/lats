package main

import (
	"github.com/jrottersman/lats/cmd"
	"github.com/jrottersman/lats/aws"
)

func main() {
	aws.Init("us-west-2")
	cmd.Execute()
}
