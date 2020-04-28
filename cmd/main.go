package main

import (
	"flag"

	gfile "github.com/ScentWoman/GFile"
)

var (
	conf = flag.String("conf", "config.json", "config file")
	addr = flag.String("addr", "127.0.0.1:8080", "listen address")
)

func init() {
	flag.Parse()
}

func main() {
	config := gfile.Parse(*conf)
	config.ListenAndServe(*addr)
}
