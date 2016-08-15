package daemon

import (
	"fmt"
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
		Scope:			"local",
		DefaultVolSz:	1024,
		Client:         client,
		Mutex:          &sync.Mutex{},
	}
	return d
}

func (d NdvolDriver) Capabilities(r volume.Request) volume.Response {
	log.Debug(DN, "Received Capabilities req")
	return volume.Response{Capabilities: volume.Capability{Scope: d.Scope}}
}

func (d NdvolDriver) Create(r volume.Request) volume.Response {
	log.Debugf(DN, "Create volume %s on %s\n", r.Name, "nedge")
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	err := d.Client.CreateVolume(r.Name, r.Options["size"])
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	return volume.Response{}
}

func (d NdvolDriver) Get(r volume.Request) volume.Response {
	log.Debug(DN, "Get volume: ", r.Name, " Options: ", r.Options)
	num, err := d.Client.GetVolume(r.Name)
	if err != nil || num < 1 {
		log.Info("Failed to retrieve volume named ", r.Name, "during Get operation: ")
		return volume.Response{}
	}
	log.Debug("Device number is: ", num)
	mnt := fmt.Sprintf("/dev/nbd%d", num)
	return volume.Response{Volume: &volume.Volume{
		Name: r.Name, Mountpoint: mnt}}
}

func (d NdvolDriver) List(r volume.Request) volume.Response {
	log.Info(DN, "List volume: ", r.Name, " Options: ", r.Options)
	vmap, err := d.Client.ListVolumes()
	if err != nil {
		log.Info("Failed to retrieve volume list", err)
		return volume.Response{Err: err.Error()}
	}
	var vols []*volume.Volume
	for name, num := range vmap {
		if name != "" {
			vols = append(vols, &volume.Volume{Name: name, Mountpoint: fmt.Sprintf("/dev/nbd%d", num)})
		}
	}
	return volume.Response{Volumes: vols}
}

func (d NdvolDriver) Mount(r volume.MountRequest) volume.Response {
	log.Info(DN, "Mount volume: ", r.Name)
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	num, err := d.Client.GetVolume(r.Name)
	if err != nil {
		log.Info("Failed to retrieve volume named ", r.Name, "during Get operation: ", err)
		return volume.Response{Err: err.Error()}
	}
	nbd := fmt.Sprintf("/dev/nbd%d", num)
	mnt, err := d.Client.MountVolume(r.Name, nbd)
	if err != nil {
		log.Info("Failed to mount volume named ", r.Name, ": ", err)
		return volume.Response{Err: err.Error()}
	}
	return volume.Response{Mountpoint: mnt}
}

func (d NdvolDriver) Path(r volume.Request) volume.Response {
	log.Info(DN, "Path volume: ", r.Name, " Options: ", r.Options)
	num, err := d.Client.GetVolume(r.Name)
	if err != nil {
		log.Info("Failed to retrieve volume named ", r.Name, "during Get operation: ", err)
		return volume.Response{Err: err.Error()}
	}
	mnt := fmt.Sprintf("/dev/nbd%d", num)
	return volume.Response{Mountpoint: mnt}
}

func (d NdvolDriver) Remove(r volume.Request) volume.Response {
	log.Info(DN, "Remove volume: ", r.Name, " Options: ", r.Options)
	d.Mutex.Lock()
	d.Client.DeleteVolume(r.Name)
	defer d.Mutex.Unlock()
	return volume.Response{}
}

func (d NdvolDriver) Unmount(r volume.UnmountRequest) volume.Response {
	log.Info(DN, "Unmount volume: ", r.Name)
	d.Mutex.Lock()
	d.Client.UnmountVolume(r.Name)
	defer d.Mutex.Unlock()

	return volume.Response{}
}

