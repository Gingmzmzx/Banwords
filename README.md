# Banwords
违禁词检查服务器，使用AC自动机算法，用Golang实现

## 使用
1. 在[Releases](https://github.com/Gingmzmzx/Banwords/releases)页面下载对应平台的二进制文件
2. 创建`config.json`文件，并写入
    ```json
    {
        "server": {
            "port": 8080,
            "host": "localhost",
            "user": "foo",
            "key": "bar"
        },
        "data": {
            "path": "data.txt"
        }
    }
    ```
3. 创建`data.txt`文件，并写入违禁词，每行一个
4. 运行二进制文件

## API
### GET /ping
检查服务器是否存活  
Return: `{"status": "ok"}`

### GET /check/:name
检查文本是否包含违禁词  
Authorization: 在header中添加`Authorization`字段，值为`Basic base64(user:key)`
Return: `{"result": true, "text": ["str1", "str2"]}`  
Return: `{"result": false}`  
例：
```shell
curl -X GET http://localhost:8080/check/Hello%20World \
    -H "Authorization: Basic aGVsbG86d29ybGQ="
```

## config.json
- `server.port`: 服务器端口
- `server.host`: 服务器地址
- `server.user`: 用户名
- `server.key`: 密码
- `data.path`: 违禁词文件路径