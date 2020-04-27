package gfile

import (
	"github.com/ScentWoman/GFile/zfile"
	"google.golang.org/api/drive/v3"
)

// Config illustrates a config
type Config struct {
	Global zfile.Config `json:"globalConfig"`
	Drive  []struct {
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
}
