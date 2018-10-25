package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type buildParams struct {
	Timestamp string
	Version   string
	Debug     bool
}

type apiParams struct {
	Reference string
	Token     string
}

type getOptParameters struct {
	Build buildParams
	API   apiParams
	Help  bool
}

// logger should be global
var logger *log.Logger

func getOptParams() *getOptParameters {
	params := &getOptParameters{}
	flag.BoolVar(&params.Build.Debug, "debug", false, "once more, with feeling")
	flag.StringVar(&params.API.Reference, "ref", "", "the reference to look up")
	flag.StringVar(&params.API.Token, "api-token", "TEST", "the token to use for the ESV api")
	flag.BoolVar(&params.Help, "help", false, "show this message")
	flag.Parse()

	if params.Help {
		fmt.Println("built:", buildTimestamp)
		fmt.Println("version:", buildVersion)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// value, ok := os.LookupEnv("")

	// set globally via linker during compilation; in version.go
	params.Build.Timestamp = getBuildTimestamp()
	params.Build.Version = getBuildVersion()

	logger = newCLILogger(params.Build.Debug)

	return params
}

// New creates a new debuglogger
func newCLILogger(echo bool) *log.Logger {
	l := log.New(ioutil.Discard, "null ", log.Lshortfile|log.LUTC|log.LstdFlags)
	if echo {
		l = log.New(os.Stderr, "debug ", log.Lshortfile|log.LUTC|log.LstdFlags)
	}
	return l
}
