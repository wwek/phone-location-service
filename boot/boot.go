package boot

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/net/ghttp"
)

// 用于应用初始化。
func init() {
	v := g.View()
	c := g.Config()
	s := g.Server()

	// 配置对象及视图对象配置
	c.AddPath("config")
	//v.AddPath("template")
	v.SetDelimiters("${", "}")

	// glog配置
	logpath := c.GetString("setting.logpath")
	glog.SetPath(logpath)
	glog.SetStdPrint(true)

	// Web Server配置
	//s.SetServerRoot("public")
	s.SetLogPath(logpath)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_ALLLOWER)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(8199)
}