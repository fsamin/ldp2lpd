package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type conf struct {
	Address    string
	TimeFactor int
	TimeUnit   string
	Verbose    bool
}

var defaultConf = conf{
	Address:    "",
	TimeFactor: 30,
	TimeUnit:   "second",
}

// Build details
var buildVersion = "dev"
var buildCommit = "unknown"
var buildDate = "unknown"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (Version %s):\n", os.Args[0], buildVersion)
		flag.PrintDefaults()
	}
}

func getConf() conf {
	address := flag.String("address", defaultConf.Address, "URI of the websocket")
	timeFactor := flag.Int("time-factor", defaultConf.TimeFactor, "time factor of interval")
	timeUnit := flag.String("time-unit", defaultConf.TimeUnit, "time unit of interval")
	debug := flag.Bool("verbose", defaultConf.Verbose, "verbose mode")

	flag.Parse()

	c := defaultConf

	// Override configuration with flags
	if *address != defaultConf.Address {
		c.Address = *address
		c.TimeFactor = *timeFactor
		c.TimeUnit = *timeUnit
		c.Verbose = *debug
	}

	if flag.NArg() > 0 {
		if flag.Arg(0) == "version" {
			fmt.Fprintf(os.Stderr, "ldp-tail version %s (%s - %s)\n", buildVersion, buildCommit, buildDate)
			os.Exit(0)
		} else if flag.Arg(0) == "help" {
			flag.Usage()
			os.Exit(0)
		} else {
			fmt.Printf("Invalid command %q\n", flag.Arg(0))
			flag.Usage()
			os.Exit(-1)
		}
	}

	if c.Address == "" {
		fmt.Fprintf(os.Stderr, "No `address` specified. Please specify it with --address or thru a config file\n")
		flag.Usage()
		os.Exit(-1)
	}

	return c
}
