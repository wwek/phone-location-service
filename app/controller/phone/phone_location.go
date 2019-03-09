package phone

import (
	"github.com/gogf/gf/g/net/ghttp"
	"phone-location-service/app/library/phone"
	"phone-location-service/app/library/response"
)

type Controller struct {

}

// 通过电话获取归属地信息
func (c *Controller) GetPhoneLocation(r *ghttp.Request) {
	if pr,err := phone.GetPhoneLocation(r.GetQueryMap()); err !=nil {
		response.Json(r,-1,err.Error())
	} else {
		response.Json(r,0,"ok",pr)
	}
}