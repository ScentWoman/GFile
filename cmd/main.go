package main

import gfile "github.com/ScentWoman/GFile"

func main() {
	config := gfile.Parse("config.json")
	config.ListenAndServe("127.0.0.1:8080")
}
