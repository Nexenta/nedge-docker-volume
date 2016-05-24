package daemon

import (
	log "github.com/Sirupsen/logrus"
	"sync"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/Nexenta/nedge-docker-volume/nedv/nedapi"
)

type NexentaDriver struct {
	TenantID       int64
	DefaultVolSz   int64
	MountPoint     string
	InitiatorIFace string
	Client         *nedapi.Client
	Mutex          *sync.Mutex
}

func DriverAlloc(cfgFile string) NexentaDriver {

	client, _ := nedapi.ClientAlloc(cfgFile)

	mntPoint := "/foo"
	initiator := "iscsiInterface"

	d := NexentaDriver{
		TenantID:       1234,
		DefaultVolSz:	1024,
		Client:         client,
		Mutex:          &sync.Mutex{},
		MountPoint:     mntPoint,
		InitiatorIFace: initiator,
	}

	return d
}

func (d NexentaDriver) Create(r volume.Request) volume.Response {
	log.Infof("Create volume %s on %s\n", r.Name, "nedge")

	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	return volume.Response{}
}


func (d NexentaDriver) Get(r volume.Request) volume.Response {
	log.Info("Get volume: ", r.Name)

	Name := "FooVol"
	mntPoint := "/barmnt"

	return volume.Response{Volume: &volume.Volume{Name: Name, Mountpoint: mntPoint}}
}

func (d NexentaDriver) List(r volume.Request) volume.Response {
	log.Info("Get volume: ", r.Name)
	var vols []*volume.Volume
	return volume.Response{Volumes: vols}
}

func (d NexentaDriver) Mount(r volume.Request) volume.Response {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	log.Infof("Mounting volume %s on %s\n", r.Name, "nedge")

	return volume.Response{Mountpoint: "/" + r.Name}
}

func (d NexentaDriver) Path(r volume.Request) volume.Response {
	log.Info("Retrieve path info for volume: ", r.Name)
	path := r.Name
	log.Debug("Path reported as: ", path)
	return volume.Response{Mountpoint: path}
}

func (d NexentaDriver) Remove(r volume.Request) volume.Response {

	log.Info("Remove/Delete Volume: ", r.Name)

	return volume.Response{}
}

func (d NexentaDriver) Unmount(r volume.Request) volume.Response {
	log.Info("Unmounting volume: ", r.Name)

	return volume.Response{}
}

