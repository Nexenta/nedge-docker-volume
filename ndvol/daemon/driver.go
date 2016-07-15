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
	DefaultVolSz	int64
	Client		*ndvolapi.Client
	Mutex		*sync.Mutex
}

func DriverAlloc(cfgFile string) NdvolDriver {

	client, _ := ndvolapi.ClientAlloc(cfgFile)
	d := NdvolDriver{
		Scope:		"local",
		DefaultVolSz:	1024,
		Client:         client,
		Mutex:          &sync.Mutex{},
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
	d.Client.CreateVolume(r.Name, d.DefaultVolSz);
	defer d.Mutex.Unlock()
	return volume.Response{}
}


func (d NdvolDriver) Get(r volume.Request) volume.Response {
	log.Info(DN, "Get volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	/*
	Name := "FooVol"
	mntPoint := "/barmnt"
	return volume.Response{Volume: &volume.Volume{Name: Name, Mountpoint: mntPoint}}
	*/
	d.Client.GetVolume(r.Name);
	return volume.Response{}
}

func (d NdvolDriver) List(r volume.Request) volume.Response {
	log.Info(DN, "List volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	var vols []*volume.Volume
	d.Client.ListVolumes()
	return volume.Response{Volumes: vols}
}

func (d NdvolDriver) Mount(r volume.Request) volume.Response {
	log.Info(DN, "Mount volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	d.Mutex.Lock()
	d.Client.MountVolume(r.Name)
	defer d.Mutex.Unlock()
	return volume.Response{}
}

func (d NdvolDriver) Path(r volume.Request) volume.Response {
	log.Info(DN, "Path volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	d.Client.GetVolume(r.Name)
	return volume.Response{}
}

func (d NdvolDriver) Remove(r volume.Request) volume.Response {
	log.Info(DN, "Remove volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	d.Mutex.Lock()
	d.Client.DeleteVolume(r.Name)
	defer d.Mutex.Unlock()
	return volume.Response{}
}

func (d NdvolDriver) Unmount(r volume.Request) volume.Response {
	log.Info(DN, "Unmount volume: ", r.Name, " MountID: ", r.MountID, " Options: ", r.Options)
	d.Mutex.Lock()
	d.Client.UnmountVolume(r.Name)
	defer d.Mutex.Unlock()

	return volume.Response{}
}

