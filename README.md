# Go ProxyPool

## 介绍

Go ProxyPool 是一个简单的代理池，它提供了一个简单的 API 接口，可以通过 API 接口获取到代理 IP，通过协程定时获取代理网站，将获取到的代理
IP 存储到指定的存储器中，并提供代理有效性检测功能，将无效的代理 IP 过滤掉。

## 项目结构

```bash
.
├── README.md
├── cmd                        # 命令行工具
├── config.yaml.sample         # 配置文件样本
├── go.mod
├── go.sum
└── pkg                        # 代码包
    ├── api                    # API服务
    ├── getters                # 获取代理的方法
    ├── global                 # 全局变量和初始化
    ├── models                 # 数据模型
    └── utils                  # 工具方法
```

## 配置文件

```yaml
storage: "redis" # 可选: memory 或 redis

redis:
  host: 127.0.0.1
  db: 1
  port: 6379

qiyun_apikey: # 齐云api密钥, 用于获取代理IP
```

## 许可

本项目采用 [MIT License](https://opensource.org/license/mit/).

## 联系

如果你有任何问题或者建议，请通过Github issue或者其他方式联系我们。

