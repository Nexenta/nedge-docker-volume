package ndvolapi

import (
	"fmt"
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
	Path		string
	ChunkSize	int64
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
	ChunkSize	int64
	Server		string
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
		Endpoint:		fmt.Sprintf("http://%s:%d/", conf.NedgeHost, conf.NedgePort),
		Path:			conf.ClusterName + "/" + conf.TenantName + "/" + conf.BucketName,
		ChunkSize:		conf.ChunkSize,
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
		log.Panic(err)
	}

	url := c.Endpoint + r["links"][0]["href"]
	resp, err := http.Get(url)
	if err != nil {
		log.Panic("Error while handling request", err)
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
		log.Panic(resp.StatusCode, string(body), err)
		return true
	}
	return false
}


func (c *Client) CreateVolume(name string, size string) (err error) {
	log.Info(fmt.Sprintf("%s: Creating volume %s", DN, name))
	data := make(map[string]interface{})
	data["volSizeMB"] = size
	data["blockSize"] = c.ChunkSize
	data["chunkSize"] = c.ChunkSize
	data["objectPath"] = c.Path + "/" + name
	_, err = c.Request("POST", fmt.Sprintf("nbd?remote=%s", c.GetRemoteAddr()), data)
	return err
}

func (c *Client) GetRemoteAddr() (addr string) {
	body, err := c.Request("GET", "system/stats", nil)
	if err != nil {
		log.Panic("Error while handling request", err)
	}
	r := make(map[string]map[string]interface{})
	jsonerr := json.Unmarshal(body, &r)
	if (jsonerr != nil) {
		log.Error(jsonerr)
	}
	stats, _ := r["response"]["stats"].(map[string]interface{})
	servers := stats["servers"].(map[string]interface{})
	for k := range servers {
		if k == c.Config.Server {
			addr = servers[k].(map[string]interface{})["ipv6addr"].(string)
		}
	}
	return addr
}

func (c *Client) GetDevNumber(name string) (number float64) {
	path := c.Path + "/" + name
	nbdList, _ := c.GetNbdList()
	for _, v := range nbdList {
		if v["objectPath"].(string) == path {
			return v["number"].(float64)
		}
	}
	return 0
}

func (c *Client) GetNbdList() (nbdList []map[string]interface{}, err error){
	body, err := c.Request("GET", fmt.Sprintf("sysconfig/nbd/devices?remote=%s", c.GetRemoteAddr()), nil)
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
	data["objectPath"] = c.Path + "/" + name
	data["number"] = c.GetDevNumber(name)
	_, err = c.Request("DELETE", fmt.Sprintf("nbd?remote=%s", c.GetRemoteAddr()), data)
	return err
}

func (c *Client) MountVolume(name string) (num int16, err error) {
	log.Debug(DN, "Mounting Volume ", name)
	/* TODO: nbd/register request */
	return num, err
}

func (c *Client) UnmountVolume(name string) (err error) {
	log.Debug(DN, "Unmounting Volume ", name)
	/* TODO: nbd/unregister request */
	return err
}

func (c *Client) GetVolume(name string) (num int16, err error) {
	log.Debug(DN, "GetVolume ", name)
	nbdList, err := c.GetNbdList()
	for _, v := range nbdList {
		if strings.Split(v["objectPath"].(string), fmt.Sprintf("%s/", c.Path))[1] == name {
			num = int16(v["number"].(float64))
			return num, err
		}
	}
	return num, err
}

func (c *Client) ListVolumes() (vmap map[string]int16, err error) {
	log.Debug(DN, "ListVolumes ")
	nbdList, err := c.GetNbdList()
	vmap = make(map[string]int16)
	for _, v := range nbdList {
		vname := strings.Split(v["objectPath"].(string), fmt.Sprintf("%s/", c.Path))[1]
		vmap[vname] = int16(v["number"].(float64))
	}
	return vmap, err
}
