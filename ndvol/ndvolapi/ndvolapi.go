package ndvolapi

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)


type Client struct {
	SVIP              string
	Endpoint          string
	DefaultAPIPort    int
	DefaultVolSize    int64 //bytes
	DefaultAccountID  int64
	DefaultTenantName string
	Config            *Config
}


type Config struct {
	IOProtocol	string // NFS, iSCSI, NBD, S3
	EndPoint	string // server:/export, IQN, devname, 
	TenantName	string
	AccessKey	string
	SecretKey	string
	MountPoint	string
}

func ReadParseConfig(fname string) (Config, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error processing config file: ", err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}
	return conf, nil
}


func ClientAlloc(configFile string) (c *Client, err error) {
	conf, err := ReadParseConfig(configFile)
	if err != nil {
		log.Fatal("Error initializing client from Config file: ", configFile, "(", err, ")")
	}

	//TODO:
	//DefaultApiPort
	//DefaultAccountID
	NdvolClient := &Client{
		SVIP:		conf.IOProtocol,
		Endpoint:	conf.EndPoint,
		DefaultAPIPort:	8888,
		DefaultAccountID:	9999,
		DefaultTenantName: conf.TenantName,
		Config:	&conf,
	}

	return NdvolClient, nil
}

