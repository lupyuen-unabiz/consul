package main

import (
	"fmt"
	//// "io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/lib"
	"github.com/mitchellh/cli"
	"strings"
	"github.com/hashicorp/consul/ipaddr"
)

func init() {
	lib.SeedMathRand()
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	//// log.SetOutput(ioutil.Discard)

	//// TODO Lup Yuen: args := os.Args[1:]
	log.Println(strings.Join(os.Environ(), "\n")) ////

	//  Get IP address.
	iplist, err := ipaddr.GetPrivateIPv4(); if err != nil { log.Panic(err) }
	log.Printf("iplist: %v", iplist)
	host := iplist[0].String()

	//  Set node name
	node := strings.Replace(os.Getenv("GCLOUD_PROJECT"), "unabiz-", "", -1)

	//  Bootstrap first agent
	cmd := "agent -bootstrap-expect=1"
	//  Subsequent agents will join
	// cmd := "join ..."

	cmdline := cmd + " " +
		"-node=" + node + " " +
		"-client=" + host + " -http-port=8080 " +
		"-server -ui -enable-script-checks=true " +
		"-data-dir=./" + node + "/data " +
		"-config-dir=./" + node + "/conf"
	log.Printf("cmdline: %s", cmdline)
	args := strings.Split(cmdline, " ")

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
