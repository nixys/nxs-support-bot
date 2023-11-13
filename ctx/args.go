package ctx

import (
	"fmt"
	"os"

	"github.com/nixys/nxs-support-bot/misc"
	"github.com/pborman/getopt/v2"
)

const (
	confPathDefault = "/nxs-support-bot.conf"
)

// Args contains arguments value read from command line
type Args struct {
	ConfigPath string
}

var version string

// ArgsRead reads arguments from command line
func ArgsRead() (Args, error) {

	var a Args

	args := getopt.New()

	helpFlag := args.BoolLong(
		"help",
		'h',
		"Show help")

	versionFlag := args.BoolLong(
		"version",
		'v',
		"Show program version")

	confPath := args.StringLong(
		"conf",
		'c',
		"",
		"Config file path")

	args.Parse(os.Args)

	/* Show help */
	if *helpFlag == true {
		argsHelp(args)
		return Args{}, misc.ErrArgSuccessExit
	}

	/* Show version */
	if *versionFlag == true {
		argsVersion()
		return Args{}, misc.ErrArgSuccessExit
	}

	/* Config path */
	if args.IsSet("conf") == true {
		a.ConfigPath = *confPath
	} else {
		a.ConfigPath = confPathDefault
	}

	return a, nil
}

func argsHelp(args *getopt.Set) {

	additionalDescription := `
	
Additional description

  Just a sample.
`

	args.PrintUsage(os.Stdout)
	fmt.Println(additionalDescription)
}

func argsVersion() {
	fmt.Println(version)
}
