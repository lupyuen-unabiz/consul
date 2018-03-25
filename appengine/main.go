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

/*
POST https://www.googleapis.com/dns/v1beta2/projects/unabiz-unaops/managedZones/unaops-zone/changes
{
  "additions": [
    {
      "name": "unaops-consul.unaops.local.",
      "type": "A",
      "ttl": 300,
      "rrdatas": [
        "1.2.3.4"
      ]}]}
 */
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

	//  Set node name.  If service name in app.yaml is consulXXX, then node name is nodeXXX
	// log.Println(strings.Join(os.Environ(), "\n"))
	node := strings.Replace(os.Getenv("GCLOUD_PROJECT"), "unabiz-", "", -1)
	service := os.Getenv("GAE_SERVICE")
	instance := strings.Replace(service, "consul", "", -1)
	if instance != "" { node = node + instance }

	//  Bootstrap the cluster using first 2 agents
	//  cmd := "agent -bootstrap-expect=2"

	//  Subsequent agents join the cluster
	cmd := "agent"

	//  Join the cluster through other agents
	para := []string{
		"-join=x.x.x.x",
		"-join=x.x.x.x",
		"-join=x.x.x.x",
		"-join=x.x.x.x",
	}

	cmdline := cmd + " " +
		"-node=" + node + " " +
		"-client=" + internalIP + " -http-port=8080 " + //  Needed for App Engine health check
		"-advertise=" + externalIP + " " +  //  Ports 8301, 8302 need to be reached through the external IP address
		"-encrypt=xxxx " +  //  Encryption key must be same for all nodes
		"-server -ui -enable-script-checks=true " +
		"-data-dir=./" + node + "/data " +
		"-config-dir=./" + node + "/conf"
	if len(para) > 0 { cmdline = cmdline + " " + strings.Join(para, " ") }
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
