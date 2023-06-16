package ctx

import (
	"fmt"
	"os"

	"github.com/pborman/getopt/v2"
)

const (
	confPathDefault = "/nxs-support-bot.conf"
)

// Args contains arguments value read from command line
type Args struct {
	ConfigPath      string
	CounterInterval *int64
}

// ArgsRead reads arguments from command line
func ArgsRead() Args {

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

	counterInterval := args.Int64Long(
		"counter-interval",
		'i',
		0,
		"User counter interval")

	args.Parse(os.Args)

	/* Show help */
	if *helpFlag == true {
		argsHelp(args)
		os.Exit(0)
	}

	/* Show version */
	if *versionFlag == true {
		argsVersion()
		os.Exit(0)
	}

	/* Config path */
	if args.IsSet("conf") == true {
		a.ConfigPath = *confPath
	} else {
		a.ConfigPath = confPathDefault
	}

	if args.IsSet("counter-interval") == true {
		a.CounterInterval = counterInterval
	}

	return a
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
	fmt.Println("1.0")
}
