package gfile

import (
	"io/ioutil"

	"log"

	"github.com/ScentWoman/GFile/zfile"
	"google.golang.org/api/drive/v3"
)

// Config illustrates a config
type Config struct {
	PageSize  int          `json:"pageSize"`
	CacheSize int          `json:"cacheSize"`
	Expire    int          `json:"expire"`
	Global    zfile.Config `json:"globalConfig"`
	Drive     []struct {
		Info         zfile.Drive `json:"info"`
		Credentials  string      `json:"credentialsFile"`
		Token        string      `json:"tokenFile"`
		DirectPrefix struct {
			Major string `json:"major"`
			Minor string `json:"minor"`
		} `json:"directPrefix"`
	} `json:"drive"`

	backend []*drive.Service
	rcache  []*cache
	direct  []*directChecker
}

// Parse parses config file.
func Parse(file string) *Config {
	var c Config
	body, e := ioutil.ReadFile(file)
	if e != nil {
		log.Fatal(e)
	}
	if e := json.Unmarshal(body, &c); e != nil {
		log.Fatal(e)
	}
	// log.Printf("%#v", c)
	return &c
}
