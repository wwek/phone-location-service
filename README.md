# phone-location-service
phone-location-service是一个电话号码归属地微服务应用

HTTP 1.1/API

Golang >= 1.12 

## 快速使用
只提供Docker容器镜像
二进制文件请自行编译

运行
```
docker run --rm -p 8199:8199  wwek/phone-location-service:latest
```

请求
```
curl "http://127.0.0.1:8199/v1/phone?phoneNumber=13888888888"
```
返回
```
{
    "data": {
        "PhoneNum": "13888888888",
        "Province": "云南",
        "City": "昆明",
        "ZipCode": "650000",
        "AreaZone": "0871",
        "CardType": "中国
        移动"
    },
    "code": 0,
    "msg": "ok"
}
```

code: 错误码(0:成功, <0:失败, >=1:错误码)

## 感谢
https://github.com/xluohome/phonedata
https://github.com/gogf/gf