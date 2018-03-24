package main

import (
	"fmt"
	//// "io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"cloud.google.com/go/compute/metadata"
	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/consul/ipaddr"
)

func init() {
	lib.SeedMathRand()
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	//// TODO Lup Yuen: Support Google Cloud AppEngine
	// Previously: log.SetOutput(ioutil.Discard)
	// Previously: args := os.Args[1:]

	//  Get Internal and External IP addresses.
	iplist, err := ipaddr.GetPrivateIPv4(); if err != nil { log.Panic(err) }
	internalIP := iplist[0].String()
	externalIP, err := metadata.ExternalIP();
	if err != nil {
		log.Printf("*** Unable to get External IP")
		externalIP = internalIP
	} else {
		log.Printf("*** External IP: %s", externalIP)
	}
	log.Printf("Internal IP: %v", iplist)

	//  Set node name
	// log.Println(strings.Join(os.Environ(), "\n"))
	node := strings.Replace(os.Getenv("GCLOUD_PROJECT"), "unabiz-", "", -1)

	//  Bootstrap first agent
	//  cmd := "agent -bootstrap-expect=1"

	//  Subsequent agents:
	cmd := "agent"
	para := ""

	//  Subsequent agents will join the first agent using this command line:
	//  	./consul join -wan -http-addr=http://Luppys-MacBook-Pro.local:8080 x.x.x.x
	//  where x.x.x.x is the External IP address obtained from App Engine -> Instances -> VM IP

	cmdline := cmd + " " +
		"-node=" + node + " " +
		"-client=" + internalIP + " -http-port=8080 " + //  Needed for App Engine health check
		"-advertise-wan=" + externalIP + " " +  //  Ports 8301, 8302 use a different external IP address
		"-server -ui -enable-script-checks=true " +
		"-data-dir=./" + node + "/data " +
		"-config-dir=./" + node + "/conf"
	if cmd == "join" {
		cmdline = cmd + " " +
			"-http-addr="
	}
	if para != "" { cmdline = cmdline + " " + para }
	log.Printf("cmdline: %s", cmdline)
	args := strings.Split(cmdline, " ")

	//// Lup Yuen: End of changes

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
