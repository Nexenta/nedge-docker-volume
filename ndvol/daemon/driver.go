package daemon

import (
	log "github.com/Sirupsen/logrus"
	"sync"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/Nexenta/nedge-docker-volume/ndvol/ndvolapi"
)

var (
	DN = "ndvoldriver "
)

type NdvolDriver struct {
	Scope		string
	TenantID	int64
	DefaultVolSz	int64
	MountPoint	string
	InitiatorIFace	string
	Client		*ndvolapi.Client
	Mutex		*sync.Mutex
}

func DriverAlloc(cfgFile string) NdvolDriver {

	client, _ := ndvolapi.ClientAlloc(cfgFile)

	mntPoint := "/foo"
	initiator := "iscsiInterface"

	d := NdvolDriver{
		Scope:		"local",
		TenantID:       1234,
		DefaultVolSz:	1024,
		Client:         client,
		Mutex:          &sync.Mutex{},
		MountPoint:     mntPoint,
		InitiatorIFace: initiator,
	}

	return d
}

func (d NdvolDriver) Capabilities(r volume.Request) volume.Response {
	log.Debug(DN, "Received Capabilities req")
	return volume.Response{}
}


func (d NdvolDriver) Create(r volume.Request) volume.Response {
	log.Infof(DN, "Create volume %s on %s\n", r.Name, "nedge")

	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	return volume.Response{}
}


func (d NdvolDriver) Get(r volume.Request) volume.Response {
	log.Info(DN, "Get volume: ", r.Name)

	Name := "FooVol"
	mntPoint := "/barmnt"

	return volume.Response{Volume: &volume.Volume{Name: Name, Mountpoint: mntPoint}}
}

func (d NdvolDriver) List(r volume.Request) volume.Response {
	log.Info(DN, "Get volume: ", r.Name)
	var vols []*volume.Volume
	return volume.Response{Volumes: vols}
}

func (d NdvolDriver) Mount(r volume.Request) volume.Response {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	log.Infof(DN, "Mounting volume %s on %s\n", r.Name, "nedge")

	return volume.Response{Mountpoint: "/" + r.Name}
}

func (d NdvolDriver) Path(r volume.Request) volume.Response {
	log.Info(DN, "Retrieve path info for volume: ", r.Name)
	path := r.Name
	log.Debug(DN, "Path reported as: ", path)
	return volume.Response{Mountpoint: path}
}

func (d NdvolDriver) Remove(r volume.Request) volume.Response {

	log.Info("DN, Remove/Delete Volume: ", r.Name)

	return volume.Response{}
}

func (d NdvolDriver) Unmount(r volume.Request) volume.Response {
	log.Info("DN, Unmounting volume: ", r.Name)

	return volume.Response{}
}

