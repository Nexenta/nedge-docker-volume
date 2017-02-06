package ndvolapi

import (
	"fmt"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"code.cloudfoundry.org/bytefmt"
	"io/ioutil"
	"errors"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"strconv"
)

const defaultSize string = "8192";
const defaultChunkSize int64 = 32768;
const defaultBlockSize int64 = 4096;
const defaultFSType string = "xfs";
const defaultMountPoint string = "/var/lib/ndvol"

var (
	DN = "ndvolapi "
)

type Client struct {
	IOProtocol	string
	Endpoint	string
	Path		string
	chunksize	int64
	blocksize	int64
	Config		*Config
}

type Config struct {
	Name		string // ndvol
	NedgeHost	string // localhost
	NedgePort	int16 // 8080
	IOProtocol	string // NFS, iSCSI, NBD, S3
	ClusterName	string
	TenantName	string
	BucketName	string
	chunksize	int64
	blocksize	int64
	Server		string
	MountPoint	string
}

func ReadParseConfig(fname string) (Config, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(DN, "Error reading config file: ", fname, " error: ", err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal(DN, "Error parsing config file: ", fname, " error: ", err)
	}
	return conf, nil
}

func ClientAlloc(configFile string) (c *Client, err error) {
	conf, err := ReadParseConfig(configFile)
	if err != nil {
		log.Fatal(DN, "Error initializing client from Config file: ", configFile, " error: ", err)
	}
	if conf.chunksize == 0 {
		conf.chunksize = defaultChunkSize
	}
	if conf.blocksize == 0 {
		conf.blocksize = defaultBlockSize
	}
	if conf.MountPoint == "" {
		conf.MountPoint = defaultMountPoint
	}
	NdvolClient := &Client{
		IOProtocol:		conf.IOProtocol,
		Endpoint:		fmt.Sprintf("http://%s:%d/", conf.NedgeHost, conf.NedgePort),
		Path:			conf.ClusterName + "/" + conf.TenantName + "/" + conf.BucketName,
		chunksize:		conf.chunksize,
		blocksize:		conf.blocksize,
		Config:			&conf,
	}

	return NdvolClient, nil
}

func (c *Client) Request(method, endpoint string, data map[string]interface{}) (body []byte, err error) {
	log.Debug("Issue request to Nexenta, endpoint: ", endpoint, " data: ", data, " method: ", method)
	if c.Endpoint == "" {
		log.Panic("Endpoint is not set, unable to issue requests")
		err = errors.New("Unable to issue json-rpc requests without specifying Endpoint")
		return nil, err
	}
	datajson, err := json.Marshal(data)
	if (err != nil) {
		log.Panic(err)
	}

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	url := c.Endpoint + endpoint
	req, err := http.NewRequest(method, url, nil)
	if len(data) != 0 {
		req, err = http.NewRequest(method, url, strings.NewReader(string(datajson)))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error while handling request", err)
		return nil, err
	}
	c.checkError(resp)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Panic(err)
	}
	return body, err
}

func (c *Client) checkError(resp *http.Response) (bool) {
	if resp.StatusCode > 399 {
		body, err := ioutil.ReadAll(resp.Body)
		log.Panic(resp.StatusCode, string(body), err)
		return true
	}
	return false
}

func (c *Client) CreateVolume(name string, options map[string]string) (err error) {
	log.Info(DN, ": Creating volume ", name)
	data := make(map[string]interface{})
	if options["size"] != "" {
		data["volSizeMB"], err = c.ConvertSize(options["size"])
		if err != nil {
			return err
		}
	} else {
		data["volSizeMB"] = defaultSize
	}

	if options["bucket"] == "" {
		data["objectPath"] = c.Path + "/" + name
	} else {
		data["objectPath"] = options["bucket"] + "/" + name
	}

	optionsObject := make(map[string]interface{})
	if options["repcount"] != "" {
		optionsObject["ccow-replication-count"] = options["repcount"] 
	}
	if options["ratelim"] != "" {
		optionsObject["ccow-iops-rate-lim"] = options["ratelim"]
	}
	data["blockSize"] = c.blocksize
	optionsObject["ccow-chunkmap-chunk-size"] = c.chunksize

	data["optionsObject"] = optionsObject

	_, err = c.Request("POST", "nbd", data)
	num, _, _ := c.GetVolume(name)

	nbd := fmt.Sprintf("/dev/nbd%d", num)
	mnt := filepath.Join(c.Config.MountPoint, name)
	if out, err := exec.Command("mkdir", "-p", mnt).CombinedOutput(); err != nil {
		log.Info("Error running mkdir command: ", err, "{", string(out), "}")
	}

	fstype := options["fstype"]
	if fstype == "" {
		fstype = defaultFSType
	}

	args := []string{"-t", fstype, nbd}
	if out, err := exec.Command("mkfs", args...).CombinedOutput(); err != nil {
		log.Error("Error running mkfs command: ", err, "{", string(out), "}")
		err = errors.New(fmt.Sprintf("%s: %s", err, out))
		return err
	}
	return err
}

func (c *Client) GetNbdList() (nbdList []map[string]interface{}, err error){
	body, err := c.Request("GET", "sysconfig/nbd/devices", nil)
	if err != nil {
		log.Panic("Error while handling request", err)
	}
	r := make(map[string]interface{})
	jsonerr := json.Unmarshal(body, &r)
	if (jsonerr != nil) {
		log.Error(jsonerr)
	}
	val := r["response"].(map[string]interface{})["value"]
	nbdList = make([]map[string]interface{}, 0)
	jsonerr = json.Unmarshal([]byte(val.(string)), &nbdList)
	if (jsonerr != nil) {
		log.Error(jsonerr)
	}
	return nbdList, err
}

func (c *Client) DeleteVolume(name string) (err error) {
	log.Debug(DN, "Deleting Volume ", name)
	data := make(map[string]interface{})
	num, path, err := c.GetVolume(name)
	data["objectPath"] = path
	data["number"] = num
	_, err = c.Request("DELETE", "nbd", data)
	mnt := filepath.Join(c.Config.MountPoint, name)
	if out, err := exec.Command("rm", "-rf", mnt).CombinedOutput(); err != nil {
		log.Info("Error running rm command: ", err, "{", string(out), "}")
	}

	return err
}

func (c *Client) MountVolume(name, nbd string) (mnt string, err error) {
	log.Debug(DN, "Mounting Volume ", name)

	mnt = filepath.Join(c.Config.MountPoint, name)
	args := []string{nbd, mnt}
	if out, err := exec.Command("mount", args...).CombinedOutput(); err != nil {
		log.Error("Error running mount command: ", err, "{", string(out), "}")
		err = errors.New(fmt.Sprintf("%s: %s", err, out))
		return mnt, err
	}
	return mnt, err
}

func (c *Client) UnmountVolume(name, nbd string) (err error) {
	log.Debug(DN, "Unmounting Volume ", name)
	if out, err := exec.Command("umount", nbd).CombinedOutput(); err != nil {
		log.Error("Error running umount command: ", err, "{", string(out), "}")
	}

	return err
}

func (c *Client) GetVolume(name string) (num int16, path string, err error) {
	log.Debug(DN, "GetVolume ", name)
	nbdList, err := c.GetNbdList()
	for _, v := range nbdList {
		path = v["objectPath"].(string)
		if strings.Split(path, "/")[len(strings.Split(path, "/")) - 1] == name {
			num = int16(v["number"].(float64))
			return num, path, err
		}
	}
	return num, path, err
}

func (c *Client) ListVolumes() (vmap map[string]string, err error) {
	log.Debug(DN, "ListVolumes ")
	nbdList, err := c.GetNbdList()
	vmap = make(map[string]string)
	for _, v := range nbdList {
		objPath := v["objectPath"].(string)

		vname := strings.Split(objPath, "/")[len(strings.Split(objPath, "/")) - 1]
		vmap[vname] = fmt.Sprintf("%s/%s", c.Config.MountPoint, vname)
	}
	log.Debug(vmap)
	return vmap, err
}

func (c *Client) ConvertSize(str_size string) (size int64, err error) {
	uSize, err := bytefmt.ToMegabytes(str_size)
	if err != nil {
		intSize, _ := strconv.Atoi(str_size)
		size = int64(intSize  / 1024 / 1024)
		err = nil
	} else {
		size = int64(uSize)
	}
	if size < 64 {
		err = errors.New("Size must have a minimum value of 64MB")
	}
	return size, err
}
