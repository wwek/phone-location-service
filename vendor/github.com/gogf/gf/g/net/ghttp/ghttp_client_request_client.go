// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// HTTP客户端请求.

package ghttp

import (
    "github.com/gogf/gf/g/text/gregex"
    "time"
    "bytes"
    "strings"
    "net/http"
    "mime/multipart"
    "os"
    "io"
    "github.com/gogf/gf/g/os/gfile"
    "errors"
    "fmt"
)

// http客户端
type Client struct {
    http.Client                   // 底层http client对象
    header      map[string]string // HEADER信息Map
    cookies     map[string]string // 自定义COOKIE
    prefix      string            // 设置请求的URL前缀
    authUser    string            // HTTP基本权限设置：名称
    authPass    string            // HTTP基本权限设置：密码
    browserMode bool              // 是否模拟浏览器模式(自动保存提交COOKIE)
}

// http客户端对象指针
func NewClient() (*Client) {
    return &Client{
        Client : http.Client {
            Transport: &http.Transport {
                DisableKeepAlives: true,
            },
        },
        header  : make(map[string]string),
        cookies : make(map[string]string),
    }
}

// 是否模拟浏览器模式(自动保存提交COOKIE)
func (c *Client) SetBrowserMode(enabled bool) {
    c.browserMode = enabled
}

// 设置HTTP Header
func (c *Client) SetHeader(key, value string) {
    c.header[key] = value
}

// 通过字符串设置HTTP Header
func (c *Client) SetHeaderRaw(header string) {
    for _, line := range strings.Split(strings.TrimSpace(header), "\n") {
        array, _ := gregex.MatchString(`^([\w\-]+):\s*(.+)`, line)
        if len(array) >= 3 {
            c.header[array[1]] = array[2]
        }
    }
}

// 设置COOKIE
func (c *Client) SetCookie(key, value string) {
    c.cookies[key] = value
}

// 使用Map设置COOKIE
func (c *Client) SetCookieMap(cookieMap map[string]string) {
    for k, v := range cookieMap {
        c.cookies[k] = v
    }
}

// 设置请求的URL前缀
func (c *Client) SetPrefix(prefix string) {
    c.prefix = prefix
}

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// 设置HTTP访问账号密码
func (c *Client) SetBasicAuth(user, pass string) {
    c.authUser = user
    c.authPass = pass
}

// GET请求
func (c *Client) Get(url string) (*ClientResponse, error) {
    return c.DoRequest("GET", url)
}

// PUT请求
func (c *Client) Put(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("PUT", url, data...)
}

// POST请求提交数据，默认使用表单方式提交数据(绝大部分场景下也是如此)。
// 如果服务端对Content-Type有要求，可使用Client对象进行请求，单独设置相关属性。
// 支持文件上传，需要字段格式为：FieldName=@file:
func (c *Client) Post(url string, data...string) (*ClientResponse, error) {
    if len(c.prefix) > 0 {
        url = c.prefix + url
    }
    param := ""
    if len(data) > 0 {
        param = data[0]
    }
    req := (*http.Request)(nil)
    if strings.Contains(param, "@file:") {
        buffer := new(bytes.Buffer)
        writer := multipart.NewWriter(buffer)
        for _, item := range strings.Split(param, "&") {
            array := strings.Split(item, "=")
            if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
                path := array[1][6:]
                if !gfile.Exists(path) {
                    return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
                }
                if file, err := writer.CreateFormFile(array[0], path); err == nil {
                    if f, err := os.Open(path); err == nil {
                        defer f.Close()
                        if _, err = io.Copy(file, f); err != nil {
                            return nil, err
                        }
                    } else {
                        return nil, err
                    }
                } else {
                    return nil, err
                }
            } else {
                writer.WriteField(array[0], array[1])
            }
        }
        writer.Close()
        if r, err := http.NewRequest("POST", url, buffer); err != nil {
            return nil, err
        } else {
            req = r
            req.Header.Set("Content-Type", writer.FormDataContentType())
        }
    } else {
        if r, err := http.NewRequest("POST", url, bytes.NewReader([]byte(param))); err != nil {
            return nil, err
        } else {
            req = r
            req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        }
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // COOKIE
    if len(c.cookies) > 0 {
        headerCookie := ""
        for k, v := range c.cookies {
            if len(headerCookie) > 0 {
                headerCookie += ";"
            }
            headerCookie += k + "=" + v
        }
        if len(headerCookie) > 0 {
            req.Header.Set("Cookie", headerCookie)
        }
    }
    // HTTP账号密码
    if len(c.authUser) > 0 {
        req.SetBasicAuth(c.authUser, c.authPass)
    }
    // 执行请求
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{
        cookies : make(map[string]string),
    }
    r.Response = resp
    return r, nil
}

// DELETE请求
func (c *Client) Delete(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("DELETE", url, data...)
}

func (c *Client) Head(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("HEAD", url, data...)
}

func (c *Client) Patch(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("PATCH", url, data...)
}

func (c *Client) Connect(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("CONNECT", url, data...)
}

func (c *Client) Options(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("OPTIONS", url, data...)
}

func (c *Client) Trace(url string, data...string) (*ClientResponse, error) {
    return c.DoRequest("TRACE", url, data...)
}

// GET请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) GetContent(url string, data...string) string {
    return c.DoRequestContent("GET", url, data...)
}

// PUT请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) PutContent(url string, data...string) string {
    return c.DoRequestContent("PUT", url, data...)
}

// POST请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) PostContent(url string, data...string) string {
    return c.DoRequestContent("POST", url, data...)
}

// DELETE请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) DeleteContent(url string, data...string) string {
    return c.DoRequestContent("DELETE", url, data...)
}

func (c *Client) HeadContent(url string, data...string) string {
    return c.DoRequestContent("HEAD", url, data...)
}

func (c *Client) PatchContent(url string, data...string) string {
    return c.DoRequestContent("PATCH", url, data...)
}

func (c *Client) ConnectContent(url string, data...string) string {
    return c.DoRequestContent("CONNECT", url, data...)
}

func (c *Client) OptionsContent(url string, data...string) string {
    return c.DoRequestContent("OPTIONS", url, data...)
}

func (c *Client) TraceContent(url string, data...string) string {
    return c.DoRequestContent("TRACE", url, data...)
}

// 请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) DoRequestContent(method string, url string, data...string) string {
    response, err := c.DoRequest(method, url, data...)
    if err != nil {
        return ""
    }
    defer response.Close()
    return string(response.ReadAll())
}

// 请求并返回response对象，该方法支持二进制提交数据
func (c *Client) DoRequest(method, url string, data...string) (*ClientResponse, error) {
    if strings.EqualFold("POST", method) {
        return c.Post(url, data...)
    }
    if len(c.prefix) > 0 {
        url = c.prefix + url
    }
    param := ""
    if len(data) > 0 {
        param = data[0]
    }
    req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader([]byte(param)))
    if err != nil {
        return nil, err
    }
    // 自定义header
    if len(c.header) > 0 {
        for k, v := range c.header {
            req.Header.Set(k, v)
        }
    }
    // COOKIE
    if len(c.cookies) > 0 {
        headerCookie := ""
        for k, v := range c.cookies {
            if len(headerCookie) > 0 {
                headerCookie += ";"
            }
            headerCookie += k + "=" + v
        }
        if len(headerCookie) > 0 {
            req.Header.Set("Cookie", headerCookie)
        }
    }
    // 执行请求
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    r := &ClientResponse{
        cookies : make(map[string]string),
    }
    r.Response = resp
    // 浏览器模式
    if c.browserMode {
        now := time.Now()
        for _, v := range r.Cookies() {
            if v.Expires.UnixNano() < now.UnixNano() {
                delete(c.cookies, v.Name)
            } else {
                c.cookies[v.Name] = v.Value
            }
        }
    }
    //fmt.Println(url, c.cookies)
    return r, nil
}




