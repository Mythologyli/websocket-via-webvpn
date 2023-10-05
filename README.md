# WebSocket via ZJU Web VPN

通过浙大 Web VPN 访问内网中的 WebSocket 服务

## 使用示例：访问内网中的 VMess 服务端

### 内网服务器配置

1. 下载 [v2ray-core](https://github.com/v2fly/v2ray-core/releases)

2. 生成 UUID

    ```bash
    $ ./v2ray uuid
    514ebdf7-8e76-a4f9-2504-6f682fea5c41
    ```

3. 修改配置文件为：

    ```json
    {
      "inbounds": [
        {
          "port": 65050,
          "protocol": "vmess",
          "settings": {
            "clients": [
              {
                "id": "514ebdf7-8e76-a4f9-2504-6f682fea5c41",
                "alterId": 0
              }
            ]
          },
          "streamSettings": {
            "network":"ws"
          }
        }
      ],
      "outbounds": [
        {
          "protocol": "freedom",
          "settings": {}
        }
      ]
    }
    ```

4. 启动 v2ray-core

    ```bash
    $ ./v2ray run
    ```

### 校外计算机配置

1. 下载 websocket-via-webvpn

2. 编写配置文件 `config.json`，内容为：

    ```json
    {
        "host": "127.0.0.1",
        "port": 8000,
        "username": "",
        "password": "",
        "websocket_host": "10.10.98.98",
        "websocket_port": 65050,
        "websocket_ssl": false,
        "websocket_path": "/"
    }
    ```

3. 启动 websocket-via-webvpn

    ```bash
    $ ./websocket-via-webvpn -config config.json
    ```

4. 此时可以通过本机的 `127.0.0.1:8080` 访问内网中的 VMess 服务端

### 可选步骤：添加到 [ZJU Rule](https://github.com/Mythologyli/ZJU-Rule)

1. 打开 [https://zjurule.xyz](https://zjurule.xyz)

2. 添加 `vmess://ws:514ebdf7-8e76-a4f9-2504-6f682fea5c41-0@127.0.0.1:8000#ZJU-WS`

3. 生成订阅链接并添加到 Clash

4. 在规则中将 ZJU 内网和 ZJU More Scholar 配置为 ZJU-WS

## 致谢

- [websocketproxy](https://github.com/pretty66/websocketproxy)
- [EasierConnect](https://github.com/lyc8503/EasierConnect)