package gfile

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/ScentWoman/GFile/zfile"
	"google.golang.org/api/drive/v3"
)

var (
	// Location of time zone.
	Location *time.Location

	driveTFormat = "2006-01-02T15:04:05.000Z"
	rTFormat     = "2006-01-02 15:04:05"
)

func init() {
	if Location == nil {
		Location, _ = time.LoadLocation("Local")
	}
}

func (c *Config) listPath(n int, path, password string) (list []zfile.File, ok bool) {
	if n >= len(c.backend) {
		return
	}
	srv := c.backend[n]
	rcache := c.rcache[n]

	list, ok = rcache.get(path, password)
	if !ok {
		return
	}
	if list == nil {
		listPath(srv, rcache, path, password)
	}

	return rcache.get(path, password)
}

func listPath(srv *drive.Service, rcache *cache, path, password string) (err error) {
	var glist *drive.FileList
	spath := strings.Split(path, "/")

	parent := "root"
	npath := "/"
	for k := range spath {
		if spath[k] == "" {
			continue
		}

		list, ok := rcache.getWithoutPass(npath)

		if !ok {
			return errors.New("wrong password")
		}

		if len(list) > 0 {
			for i := range list {
				if list[i].Name == spath[k] {
					parent = *list[i].URL
					break
				}
			}
		} else {
			glist, err = srv.Files.List().Fields("*").Q(`'` + parent + `' in parents`).Do()
			if err != nil {
				log.Println(err)
				return
			}
			// log.Println(glist.Files)

			var nlist []zfile.File
			var npass string
			for _, v := range glist.Files {
				nlist = append(nlist, zfile.File{
					Name: v.Name,
					Path: npath,
					Size: v.Size,
					Time: toTime(v.ModifiedTime),
					Type: toType(v.MimeType),
					URL:  &v.Id,
				})
				if strings.ToLower(v.Name) == "password.txt" {
					resp, err := srv.Files.Get(v.Id).Download()
					if err != nil {
						return err
					}
					body, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						return err
					}
					npass = string(body)
				}
				if v.Name == spath[k] {
					parent = v.Id
				}
			}
			rcache.set(npath, npass, nlist)
		}
		npath = npath + spath[k] + "/"
	}

	glist, err = srv.Files.List().Fields("*").Q(`'` + parent + `' in parents`).Do()
	if err != nil {
		log.Println(err)
		return
	}
	// log.Println(glist.Files)

	var nlist []zfile.File
	var npass string
	for _, v := range glist.Files {
		nlist = append(nlist, zfile.File{
			Name: v.Name,
			Path: npath,
			Size: v.Size,
			Time: toTime(v.ModifiedTime),
			Type: toType(v.MimeType),
			URL:  &v.Id,
		})
		if strings.ToLower(v.Name) == "password.txt" {
			resp, err := srv.Files.Get(v.Id).Download()
			if err != nil {
				return err
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}
			npass = string(body)
		}
	}
	rcache.set(path, npass, nlist)

	return
}

func toType(s string) string {
	if s == "application/vnd.google-apps.folder" {
		return "FOLDER"
	}
	return "FILE"
}

func toTime(s string) string {
	t, e := time.ParseInLocation(driveTFormat, s, Location)
	if e != nil {
		return s
	}
	return t.Format(rTFormat)
}
