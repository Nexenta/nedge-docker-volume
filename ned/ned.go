package main

import (
	"github.com/nacharya/ned/nedapi"
	"github.com/nacharya/ned/nedcli"
	"os"
)

const (
	VERSION = "0.0.1"
)

var (
	client *nedapi.Client
)

func main() {
	ncli := nedcli.NewCli(VERSION)
	ncli.Run(os.Args)
}

