package main

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli"
	"github.com/watermint/dcfg/dispatch"
	"github.com/watermint/dcfg/explorer"
	"os"
)

var (
	AppVersion string
)

func main() {
	options := cli.Options{}
	options.Parse(os.Args[:1])
	if err := options.Validate(); err != nil {
		fmt.Errorf("Error: %v", err)
		options.Usage()
		os.Exit(1)
	}

	explorer.Start(options, AppVersion)

	defer explorer.Report()
	defer seelog.Flush()

	context := cli.ExecutionContext{
		Options: options,
	}

	dispatch.Dispatch(context)
}
