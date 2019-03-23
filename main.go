package main

import (
	"github.com/gogf/gf/g"
	_ "phone-location-service/boot"
	_ "phone-location-service/router"
)

func main() {
	g.Server().SetServerAgent("nginx")
	g.Server().Run()
}
