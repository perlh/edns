# edns


- 一个简单的DNS服务器，提供一些接口 ，下次有空写一下doc。

- 这个项目使用了Gorm库对数据库进行操作，是用来redis来保存dns信息。

- [客户端代码: client.py](client.py)
- [服务器代码: main.go,...](main.go)

# API

## POST register

POST /register

> 返回示例

```json
{
  "code": 200,
  "msg": "用户注册成功",
  "Id": 0,
  "Email": "test@ddns.cool",
  "Passwd": "test",
  "Role": 0,
  "CreateTime": 0
}
```


## POST get_dns

POST /get_dns

> 返回示例

```json
{
  "code": 200,
  "msg": "",
  "data": [
    {
      "Id": 1,
      "Domain": "ddns.cool",
      "HostRecode": "zdx",
      "RecodeType": "a",
      "RecodeValue": "172.24.1.100",
      "TTL": 600,
      "LastOptionTime": 1671103055,
      "UserID": 3
    },
    {
      "Id": 2,
      "Domain": "ddns.cool",
      "HostRecode": "ddf",
      "RecodeType": "a",
      "RecodeValue": "222.204.52.222",
      "TTL": 600,
      "LastOptionTime": 1671106094,
      "UserID": 4
    }
  ]
}
```



## POST get_user_info

POST /get_user_info

> 返回示例

```json
{
  "code": 200,
  "msg": "",
  "data": [
    {
      "Id": 1,
      "Email": "905008234@qq.com",
      "Passwd": "qwert@123",
      "Role": 1,
      "CreateTime": 1671097726
    },
    {
      "Id": 2,
      "Email": "test23@qq.com",
      "Passwd": "test",
      "Role": 0,
      "CreateTime": 1671098402
    }
  ]
}
```



## POST user_delete

POST /user_delete

> Body 请求参数

```yaml
root_token: dsafasdf
email: test@qq.com

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» root_token|body|string| 否 |none|
|» email|body|string| 否 |none|

> 返回示例

```json
{
  "code": 200,
  "msg": "删除用户成功"
}
```

## POST register_dns

POST /register_dns

> Body 请求参数

```yaml
token: xxx
host_recode: test
domain: hsm.cool
recode_type: a
ttl: "50"
recode_value: 222.204.102.134
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» token|body|string| 否 |none|
|» host_recode|body|string| 否 |none|
|» domain|body|string| 否 |none|
|» recode_type|body|string| 否 |none|
|» ttl|body|string| 否 |none|
|» recode_value|body|string| 否 |none|

> 返回示例

```json
{
  "code": 200,
  "msg": "添加成功",
  "data": {
    "Id": 0,
    "Domain": "hsm.cool",
    "HostRecode": "dns",
    "RecodeType": "a",
    "RecodeValue": "222.203.22.2",
    "TTL": 50,
    "LastOptionTime": 1671093429,
    "UserID": 2
  }
}
```



## POST dns_delete

POST /dns_delete

> Body 请求参数

```yaml
token: fdsfas
dns_id: "4"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» token|body|string| 否 |none|
|» dns_id|body|string| 否 |none|

> 返回示例

```json
{
  "code": 200,
  "msg": "删除成功",
  "data": {
    "Id": 0,
    "Domain": "",
    "HostRecode": "",
    "RecodeType": "",
    "RecodeValue": "",
    "TTL": 0,
    "LastOptionTime": 0,
    "UserID": 0
  }
}
```
