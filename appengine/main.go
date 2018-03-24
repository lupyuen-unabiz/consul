package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/lib"
	"github.com/mitchellh/cli"
	"strings"
)

func init() {
	lib.SeedMathRand()
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	//// TODO Lup Yuen: args := os.Args[1:]
	log.Println(strings.Join(os.Environ(), "\n")) ////
	node := "unaops"
	cmd := "-bootstrap-expect=1"
	// cmd := "join unaops"

	args := strings.Split(
		"agent -node=" + node + " " +
			cmd + " " +
			"-server -ui -http-port 8080 -enable-script-checks=true " +
			"-data-dir=./" + node + "/data " +
			"-config-dir=./" + node + "/conf",
		" ")

	for _, arg := range args {
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "--version" {
			args = []string{"version"}
			break
		}
	}

	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}
	cmds := command.Map(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}

	cli := &cli.CLI{
		Args:         args,
		Commands:     cmds,
		Autocomplete: true,
		Name:         "consul",
		HelpFunc:     cli.FilteredHelpFunc(names, cli.BasicHelpFunc("consul")),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
