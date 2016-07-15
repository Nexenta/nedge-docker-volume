package ndvolapi

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"errors"
	"net/http"
	"strings"
	"time"
)

var (
	DN = "ndvolapi "
)

type Client struct {
	IOProtocol	string
	Endpoint	string
	ClusterName	string
	TenantName	string
	BucketName	string
	ChunkSize	int64
	Config		*Config
}

type Config struct {
	Name		string // ndvol
	NedgeHost	string // localhost
	NedgePort	int64 // 8080
	IOProtocol	string // NFS, iSCSI, NBD, S3
	ClusterName	string
	TenantName	string
	BucketName	string
	ChunkSize	int64
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
	NdvolClient := &Client{
		IOProtocol:		conf.IOProtocol,
		Endpoint:		"http://" + conf.NedgeHost + "/" + string(conf.NedgePort) + "/",
		ClusterName:		conf.ClusterName,
		TenantName:		conf.TenantName,
		BucketName:		conf.BucketName,
		ChunkSize:		conf.ChunkSize,
		Config:			&conf,
	}

	return NdvolClient, nil
}

func (c *Client) Request(method, endpoint string, data map[string]interface{}) (body []byte, err error) {
	log.Debug("Issue request to Nexenta, endpoint: ", endpoint, " data: ", data, " method: ", method)
	if c.Endpoint == "" {
		log.Error("Endpoint is not set, unable to issue requests")
		err = errors.New("Unable to issue json-rpc requests without specifying Endpoint")
		return nil, err
	}
	datajson, err := json.Marshal(data)
	if (err != nil) {
		log.Error(err)
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
		log.Error("Error while handling request", err)
		return nil, err
	}
	c.checkError(resp)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if (err != nil) {
		log.Error(err)
	}
	if (resp.StatusCode == 202) {
		body, err = c.resend202(body)
	}
	return body, err
}

func (c *Client) resend202(body []byte) ([]byte, error) {
	time.Sleep(1000 * time.Millisecond)
	r := make(map[string][]map[string]string)
	err := json.Unmarshal(body, &r)
	if (err != nil) {
		log.Error(err)
	}

	url := c.Endpoint + r["links"][0]["href"]
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Error while handling request", err)
		return nil, err
	}
	defer resp.Body.Close()
	c.checkError(resp)

	if resp.StatusCode == 202 {
		body, err = c.resend202(body)
	}
	body, err = ioutil.ReadAll(resp.Body)
	return body, err
}

func (c *Client) checkError(resp *http.Response) (bool) {
	if resp.StatusCode > 399 {
		body, err := ioutil.ReadAll(resp.Body)
		log.Error(resp.StatusCode, string(body), err)
		return true
	}
	return false
}


func (c *Client) CreateVolume(name string, size int64) (err error) {
	log.Debug(DN, "Creating volume %s", name)
	path := c.ClusterName + "/" + c.TenantName + "/" + c.BucketName + "/" + name
	data := map[string]interface{} {
		"path": path,
	}
	data = make(map[string]interface{})
	data["number"] = len(name)
	data["volSizeMB"] = size >> 20
	data["blockSize"] = c.ChunkSize
	data["chunkSize"] = c.ChunkSize
	data["objectPath"] = path
	// _, err = c.Request("POST", c.EndPoint, data)
	err = nil
	return err
}

func (c *Client) DeleteVolume(name string) (err error) {
	log.Debug(DN, "Deleting Volume ", name)
	path := c.ClusterName + "/" + c.TenantName + "/" + c.BucketName + "/" + name
	data := map[string]interface{} {
		"path": path,
	}
	data = make(map[string]interface{})
	data["number"] = len(name)
	data["volSizeMB"] = c.ChunkSize
	data["blockSize"] = c.ChunkSize
	data["chunkSize"] = c.ChunkSize
	data["objectPath"] = path
	// _, err = c.Request("DELETE", c.EndPoint, data)
	return err
}

func (c *Client) MountVolume(name string) (err error) {
	log.Debug(DN, "Mounting Volume ", name)
	err = nil
	return err
}

func (c *Client) UnmountVolume(name string) (err error) {
	log.Debug(DN, "Unmounting Volume ", name)
	err = nil
	return err
}

func (c *Client) GetVolume(name string) (vname string, err error) {
	log.Debug(DN, "GetVolume ", name)
	vname = ""
	err = nil
	return vname, err
}

func (c *Client) ListVolumes() (vlist []string, err error) {
	log.Debug(DN, "ListVolumes ")
	err = nil
	return vlist, err
}
