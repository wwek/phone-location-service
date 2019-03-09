package main

import (
	_ "phone-location-service/boot"
	_ "phone-location-service/router"
	"github.com/gogf/gf/g"
)

func main() {
	g.Server().Run()
}