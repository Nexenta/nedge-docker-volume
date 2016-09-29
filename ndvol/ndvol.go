package main

import (
	"github.com/Nexenta/nedge-docker-volume/ndvol/ndvolapi"
	"github.com/Nexenta/nedge-docker-volume/ndvol/ndvolcli"
	"os"
)

const (
	VERSION = "0.0.1"
)

var (
	client *ndvolapi.Client
)

func main() {
	ncli := ndvolcli.NewCli(VERSION)
	ncli.Run(os.Args)
}
