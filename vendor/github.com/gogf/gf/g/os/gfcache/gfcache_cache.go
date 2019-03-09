// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 文件缓存.
package gfcache

import (
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gfsnotify"
)

// 设置容量大小(byte)
func (c *Cache) SetCap(cap int) {
    c.cap.Set(cap)
}

// 获得缓存容量大小(byte)
func (c *Cache) GetCap() int {
    return c.cap.Val()
}

// 获得已缓存的文件大小(byte)
func (c *Cache) GetSize() int {
    return c.size.Val()
}

// 获得文件内容 string
func (c *Cache) GetContents(path string) string {
    return string(c.GetBinContents(path))
}

// 获得文件内容 []byte
func (c *Cache) GetBinContents(path string) []byte {
    if v := c.cache.Get(path); v != nil {
        return v.([]byte)
    }
    b := gfile.GetBinContents(path)
    // 读取到内容，并且没有超过缓存容量限制时才会执行缓存
    if len(b) > 0 && (c.cap.Val() == 0 || c.size.Val() < c.cap.Val()) {
        c.size.Add(len(b))
        c.addMonitor(path)
        c.cache.Set(path, b)
    }
    return b
}

// 添加文件监控，一旦文件有变化立即清除缓存，下一次读取的时候再执行缓存。
func (c *Cache) addMonitor(path string) {
    // 防止多goroutine同时调用
    if c.cache.Contains(path) {
        return
    }
    gfsnotify.Add(path, func(event *gfsnotify.Event) {
        if r := c.cache.Get(path); r != nil {
            c.cache.Remove(path)
            c.size.Add(-len(r.([]byte)))
        }
    })
}