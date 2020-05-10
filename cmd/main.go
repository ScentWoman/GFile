package main

import (
	"flag"
	"time"

	gfile "github.com/ScentWoman/GFile"
)

var (
	conf     = flag.String("conf", "config.json", "config file")
	addr     = flag.String("addr", "127.0.0.1:8080", "listen address")
	timezone = flag.String("timezone", "Asia/Shanghai", "timezone")
)

func init() {
	flag.Parse()
}

func main() {
	gfile.Location, _ = time.LoadLocation(*timezone)
	config := gfile.Parse(*conf)
	config.ListenAndServe(*addr)
}
