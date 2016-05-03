package nedcli

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var (
	FooCmd =  cli.Command{
		Name:  "FooCmd",
		Usage: "Foo related commands",
		Subcommands: []cli.Command{
			FooCreateCmd,
			FooDeleteCmd,
		},
	}

	FooCreateCmd = cli.Command{
		Name:  "create",
		Usage: "create a new volume: `create [options] NAME`",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "size",
				Usage: "size of volume in bytes ",
			},
			cli.StringFlag{
				Name:  "account",
				Usage: "account id to assign volume",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "Specify a volume type ",
			},
		},
		Action: cmdCreate,
	}
	FooDeleteCmd = cli.Command{
		Name:  "delete",
		Usage: "delete an existing volume: `delete NAME`",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "range",
				Value: "",
				Usage: ": deletes a range of volume`",
			},
		},
		Action: cmdDelete,
	}
)


func cmdCreate(c *cli.Context) {
	fmt.Println("cmdCreate: ", c.String("size"));
}

func cmdDelete(c *cli.Context) {
	fmt.Println("cmdCreate: ", c.String("name"));
}

