package nedcli

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var (
	VolumeCmd =  cli.Command{
		Name:  "VolumeCmd",
		Usage: "Volume related commands",
		Subcommands: []cli.Command{
			VolumeCreateCmd,
			VolumeDeleteCmd,
			VolumeListCmd,
		},
	}

	VolumeCreateCmd = cli.Command{
		Name:  "create",
		Usage: "create a new volume: `create [options] NAME`",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "size",
				Usage: "size of volume in bytes ",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "Specify a volume type ",
			},
		},
		Action: cmdCreateVolume,
	}
	VolumeDeleteCmd = cli.Command{
		Name:  "delete",
		Usage: "delete an existing volume: `delete NAME`",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "range",
				Value: "",
				Usage: ": deletes a range of volume`",
			},
		},
		Action: cmdDeleteVolume,
	}
	VolumeListCmd = cli.Command{
		Name:  "list",
		Usage: "list existing volumes",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "range",
				Value: "",
				Usage: ": range of volume`",
			},
		},
		Action: cmdListVolume,
	}

)


func cmdCreateVolume(c *cli.Context) {
	fmt.Println("cmdCreate: ", c.String("size"));
}

func cmdDeleteVolume(c *cli.Context) {
	fmt.Println("cmdDelete: ", c.String("name"));
}

func cmdListVolume(c *cli.Context) {
	fmt.Println("cmdDelete: ", c.String("name"));
}
