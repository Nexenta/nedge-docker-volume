package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"path/filepath"
)

var (
	defaultDir = filepath.Join(volume.DefaultDockerRootDirectory, "nedge")
)

func Start(cfgFile string, debug bool) {
	if debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	d := DriverAlloc(cfgFile)
	h := volume.NewHandler(d)
	log.Info(h.ServeUnix("root", "toor"))
}