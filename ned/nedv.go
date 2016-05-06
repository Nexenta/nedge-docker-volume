package main

import (
	"github.com/Nexenta/nedge-docker-volume/ned/nedapi"
	"github.com/Nexenta/nedge-docker-volume/ned/nedcli"
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

