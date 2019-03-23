package ctl_phone

import (
	"github.com/gogf/gf/g/net/ghttp"
	"phone-location-service/app/library/lib_phone"
	"phone-location-service/app/library/lib_response"
)

type Controller struct {
}

// 通过电话获取归属地信息
func (c *Controller) GetPhoneLocation(r *ghttp.Request) {
	if pr, err := lib_phone.GetPhoneLocation(r.GetQueryMap()); err != nil {
		lib_response.Json(r, -1, err.Error())
	} else {
		lib_response.Json(r, 0, "ok", pr)
	}
}
