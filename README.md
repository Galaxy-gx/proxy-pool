# proxy-pool
爬虫代理池

### 体验地址
* [get](http://81.68.131.249:9001/proxy/get)
* [getall](http://81.68.131.249:9001/proxy/getall)


### 新建数据库
* 使用mysql数据库
* 新建数据库 `CREATE DATABASE proxy_pool CHARACTER SET utf8mb4`
* 使用 schema.sql 创建表结构


### docker-compose 启动服务
```yaml
version: "3"

services:
  proxy-pool-api:
    image: registry.cn-shanghai.aliyuncs.com/release-lib/proxy-pool:latest
    container_name: proxy-pool-api
    restart: always
    ports:
      - 9001:9001
    environment:
      - TZ=Asia/Shanghai
      - CONFIG_FILE=/etc/conf.yaml
      - HOST=****
      - PORT=3306
      - USERNAME=****
      - PASSWORD=****
      - DATABASE=proxy_pool
    command: api

  proxy-pool-schduler:
    image: registry.cn-shanghai.aliyuncs.com/release-lib/proxy-pool:latest
    container_name: proxy-pool-scheduler
    restart: always
    environment:
      - TZ=Asia/Shanghai
      - CONFIG_FILE=/etc/conf.yaml
      - HOST=****
      - PORT=3306
      - USERNAME=****
      - PASSWORD=****
      - DATABASE=proxy_pool
    command: scheduler
```

### 等待几分钟访问接口
* 获取一个proxy `/proxy/get`
* 获取所有可用proxy `/proxy/getall`

### TODO
- [ ] 增加更多免费代理网址