package gfile

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/ScentWoman/GFile/zfile"
)

type directChecker struct {
	major, minor string
	value        atomic.Value
}

func newChecker(major, minor string) (d *directChecker) {
	d = &directChecker{
		major: major,
		minor: minor,
	}
	d.value.Store(major)
	return d
}

func (d *directChecker) autoCheck() {
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		resp, e := client.Get(d.major)
		if e != nil {
			d.value.Store(d.minor)
		} else {
			resp.Body.Close()
			d.value.Store(d.major)
		}
		time.Sleep(time.Minute)
	}
}

func (d *directChecker) get() string {
	return d.value.Load().(string)
}

func (d *directChecker) check(list []zfile.File) {
	prefix := d.get()
	for k := range list {
		if list[k].Type == "FILE" {
			URL := prefix + "/" + *list[k].URL + "/" + list[k].Name
			list[k].URL = &URL
		} else {
			list[k].URL = nil
		}
	}
}
