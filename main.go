package main

import (
	"github.com/jrottersman/lats/aws"
	"github.com/jrottersman/lats/cmd"
)

func main() {
	aws.Init("us-west-2")
	cmd.Execute()
}
