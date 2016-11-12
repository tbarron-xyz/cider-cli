package main

import (
	"flag"
)

var _port = flag.Int("port", 6969, "port")
var _url = flag.String("url", "localhost", "url")

func init() {
	flag.Parse()
}
