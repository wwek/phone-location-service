package router

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"phone-location-service/app/controller/phone"
	"phone-location-service/app/library/response"
)

// 统一路由注册.
func init() {

	s := g.Server()

	s.BindObjectMethod("/v1/phone",new(phone.Controller),"GetPhoneLocation")

	s.BindHandler("/*any", func(r *ghttp.Request){
		response.Json(r,-1,"404")
	})
}