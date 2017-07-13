package main

import (
	"github.com/zanecloud/metad/cli"
)

var (
	Version   string
	GitCommit string
	BuildTime string
)

func main() {

	cli.Run()
}
