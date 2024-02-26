# 服务架构增强

## Zlog 配置

日志服务配置，在 zaplog 的基础上进行封装，使用了 zaplog 的 Sugar() 语法糖，并且把方法都出来了。

### 支持的配置

```toml
[logger]
encoder = "json"
[logger.sls]
accessKeyId = "xxx"
accessKeySecret = "xxx"
endpoint = "cn-shanghai.log.aliyuncs.com"
logStoreName = "w303b"
projectName = "w303b-test"
sourceIp = "0.0.0.0"
topic = "servers"
```

#### encoder

这里的 encoder 如果是 json，那么需要填入 sls 信息。
json 格式的数据通常会上传到解析平台，现在只支持 aliyun 的 sls 服务。

#### accessKeyId & accessKeySecret

aliyun sls 服务的的认证信息

#### endpoint

节点, 可以通过连接的节点不同区分日志存储在哪里，如我需要国内就用国内的节点。

#### logStoreName & projectName

储存的 storeName，要对应上 aliyun sls 上的配置

#### sourceIp

自己的 IP

#### topic

对应 \_\_topic 索引，一般来说不中要，但是这里用来区分是 http 服务还是 websocket 服务

## Mysql 配置

将 gorm 封装成 dbx 包，里面使用了 gorm 的连接，在创建连接是如果启用了 ssh 配置，会优先使用跳板机的方式去连接远程的 mysql，这功能在有些情况下很有用，就比如我需要在本地连接测试环境的 mysql。

### 支持的配置

```toml
[mysql]
dsn = "root:qq123123@tcp(127.0.0.1:3306)/sikey?charset=utf8mb4&parseTime=true&loc=Local"
ssh = false
# Ignore ErrRecordNotFound error for logger
skipDefaultTransaction = true
# Slow SQL threshold
slowThreshold = 600
# Ignore ErrRecordNotFound error for logger
ignoreRecordNotFoundError = true
# 设置空闲连接池中连接的最大数量
maxIdleConns = 10
# 设置连接的有效时长 当 <= 0 时，连接永久保存，默认值时 0 。如果设置了 maxLifetime 会开启连接自动清理，
# 清理的代码在 connectionCleaner 中， 它开启一个定时器，定时检查空闲连接池中的连接，超期的关闭连接。
maxLifetime = -1
# 设置打开数据库连接的最大数量。
maxOpenConns = 100
```

### Redis 配置

### 支持的配置

```toml
[redis]
addr = "106.75.230.4:6379"
channel = "message"
connectKey = "connects"
db = 0
password = "xxx"
```
