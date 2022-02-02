package main

import (
	"experiment_lwc/commons"
	"experiment_lwc/config"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var configuration config.Configuration

var runOpts struct {
	ConfigPath string `long:"config" description:"Config path" required:"true"`
}

var getOpts struct {
	ID string `long:"id" description:"Container ID" required:"true"`
}

var cleanupOpts struct {
	ID string `long:"id" description:"Container ID" required:"true"`
}

var networkOpts struct {
	Action  string `long:"action" choice:"create" choice:"delete" choice:"list" description:"action" required:"true"`
	Name    string `long:"name" description:"Bridge name"`
	Address string `long:"address" description:"Bridge address"`
}

var choiceOpts struct {
	Opt string `long:"opt" choice:"list" choice:"inspect" choice:"cleanup" choice:"network" choice:"run" description:"choice" required:"true"`
}

func ParseArgs() {
	_, err := flags.NewParser(&choiceOpts, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown).ParseArgs(os.Args)
	switch choiceOpts.Opt {
	case "list":
		break
	case "inspect":
		ParseInspectArgs()
		break
	case "cleanup":
		ParseCleanupArgs()
		break
	case "run":
		ParseRunArgs()
		break
	case "network":
		ParseNetworkArgs()
		break
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func ParseRunArgs() {
	parser := flags.NewParser(&runOpts, flags.Default|flags.IgnoreUnknown)
	parser.Usage = "--opt run [OPTIONS]"
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(0)
	}
	configuration = *config.NewConfiguration()
	if runOpts.ConfigPath != "" {
		content, err := ioutil.ReadFile(runOpts.ConfigPath)
		commons.Must(err)
		commons.Must(yaml.Unmarshal(content, &configuration))
	} else {
		panic("config path is not specified")
	}
	if configuration.ID == "" {
		configuration.ID = commons.StringRandom(64, commons.Lowercase+commons.Numeric)
	}
	if configuration.Hostname == "" {
		configuration.Hostname = configuration.ID[:12]
	}
}

func ParseInspectArgs() {
	parser := flags.NewParser(&getOpts, flags.Default|flags.IgnoreUnknown)
	parser.Usage = "--opt inspect [OPTIONS]"
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(0)
	}
	configuration = *config.NewConfiguration()
	configuration.ID = getOpts.ID
}

func ParseCleanupArgs() {
	parser := flags.NewParser(&cleanupOpts, flags.Default|flags.IgnoreUnknown)
	parser.Usage = "--opt run [OPTIONS]"
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(0)
	}
	configuration = *config.NewConfiguration()
	configuration.ID = cleanupOpts.ID
}

func ParseNetworkArgs() {
	parser := flags.NewParser(&networkOpts, flags.Default|flags.IgnoreUnknown)
	parser.Usage = "--opt network [OPTIONS]"
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(0)
	}
}
