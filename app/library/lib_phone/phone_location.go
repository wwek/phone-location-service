package lib_phone

import (
	"errors"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gvalid"
	"github.com/xluohome/phonedata"
)

func GetPhoneLocation(data g.MapStrStr) (*phonedata.PhoneRecord, error) {
	// 参数效验
	rules := []string{
		"phoneNumber @required|length:7,11#电话号码不能为空|号码长度应该在:min到:max之间",
	}
	if e := gvalid.CheckMap(data, rules); e != nil {
		return nil, errors.New(e.String())
	}

	pr, err := phonedata.Find(data["phoneNumber"])
	if err != nil {
		return nil, err
	}

	return pr, err
}
