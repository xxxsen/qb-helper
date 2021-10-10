qb-helper
===

[TOC]

qb助手, 用于qbittorrent自动化处理.
目前主要提供的功能是ban客户端.

## 配置

### 模板:

```json
{
    "auth": {
        "username": "xxx",
        "password": "xxx",
        "host": "https://qb.ccc.com"
    },
    "log": {
        "console": true
    },
    "cron_config": [
        {
            "name": "anti-leecher",
            "args": {
                "ban_client": [
                    {
                        "key": "-XL0012",
                        "mode": "prefix"
                    },
                    {
                        "key": "Xunlei",
                        "mode": "prefix"
                    }
                ]
            },
            "enable": true
        }
    ]
}
```

