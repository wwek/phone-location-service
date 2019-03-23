package router

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"phone-location-service/app/controller/ctl_phone"
	"phone-location-service/app/library/lib_response"
)

// 统一路由注册.
func init() {

	s := g.Server()

	s.BindObjectMethod("/v1/phone", new(ctl_phone.Controller), "GetPhoneLocation")

	s.BindHandler("/*any", func(r *ghttp.Request) {
		lib_response.Json(r, -1, "404")
	})
}
