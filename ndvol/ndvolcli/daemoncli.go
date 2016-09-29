package ndvolcli

import (
	"github.com/urfave/cli"
	ndvolDaemon "github.com/Nexenta/nedge-docker-volume/ndvol/daemon"
	"github.com/sevlyar/go-daemon"
	"os"
	log "github.com/Sirupsen/logrus"
	// "log"
	"fmt"
	"time"
	"syscall"
	// "flag"
)

var (
	DaemonCmd = cli.Command{
		Name:  "daemon",
		Usage: "daemon related commands",
		Subcommands: []cli.Command{
			DaemonStartCmd,
			DaemonStopCmd,
		},
	}

	DaemonStartCmd = cli.Command{
		Name:  "start",
		Usage: "Start the Nedge Docker Daemon: `start [options] NAME`",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "Enable verbose/debug logging: `[--verbose]`",
			},
			cli.StringFlag{
				Name:  "config, c",
				Usage: "Config file for daemon (default: /etc/nvd/nvd.json): `[--config /etc/nvd/nvd.json]`",
			},
		},
		Action: cmdDaemonStart,
	}
	DaemonStopCmd = cli.Command{
		Name:  "stop",
		Usage: "Stop the Nedge Docker Daemon",
		Action: cmdDaemonStop,
	}
)

func cmdDaemonStop(c *cli.Context) {
	termHandler(syscall.SIGQUIT)
}

func cmdDaemonStart(c *cli.Context) {
	fmt.Println("daemon start")
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		// WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[ndvol daemon]"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(d, err)
	if d != nil {
		fmt.Println("return")
		return
	}
	defer cntxt.Release()

	fmt.Println("- - - - - - - - - - - - - - -")
	fmt.Println("daemon started")
	go worker(c)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker(c *cli.Context) {
	fmt.Println("worker")
	DaemonStart(c)
	for {
		time.Sleep(time.Second)
		if _, ok := <-stop; ok {
			break
		}
	}
	done <- struct{}{}
}

func DaemonStart(c *cli.Context) {
	verbose := c.Bool("verbose")
	cfg := c.String("config")
	if cfg == "" {
		cfg = "/opt/nedge/etc/ccow/ndvol.json"
	}
	ndvolDaemon.Start(cfg, verbose)
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	log.Println("stop")
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
