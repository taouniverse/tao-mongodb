# github.com/taouniverse/tao-mongodb

[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao-mongodb)](https://goreportcard.com/report/github.com/taouniverse/tao-mongodb)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao-mongodb?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao-mongodb?tab=doc)

Tao Universe 组件单元（Unit），基于泛型工厂模式封装 **MongoDB** 数据库。

## 安装

```bash
go get github.com/taouniverse/tao-mongodb
```

## 使用

### 导入

```go
import _ "github.com/taouniverse/tao-mongodb"
```

### 配置

```yaml
# 单实例配置
mongodb:
  host: localhost
  port: 27017
  user: tao
  password: 123456qwe
  db: test

# 多实例配置（如主从分离）
mongodb:
  default_instance: master
  master:
    host: localhost
    port: 27017
    user: tao
    password: 123456qwe
    db: mydb
  replica:
    host: backup.example.com
    port: 27017
    user: readonly
    password: ro_pass
    db: mydb
```

### 配置字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `host` | string | `localhost` | MongoDB 服务器地址 |
| `port` | int | `27017` | MongoDB 端口 |
| `user` | string | - | 用户名 |
| `password` | string | - | 密码 |
| `db` | string | - | 默认数据库 |
| `auth_source` | string | `admin` | 认证数据库 |
| `timeout` | duration | `10s` | 连接超时时间 |
| `max_pool_size` | int | `100` | 连接池最大连接数 |
| `min_pool_size` | int | `10` | 连接池最小连接数 |

## 工厂模式 API

| API | 说明 |
|-----|------|
| `mongodb.M` | 配置实例 `*Config` |
| `mongodb.Factory` | `*tao.BaseFactory[*mongo.Client]` 工厂实例 |
| `mongodb.Client()` | 获取默认 MongoDB 客户端 `(*mongo.Client, error)` |
| `mongodb.GetClient(name)` | 获取指定名称的客户端 `(*mongo.Client, error)` |

## 使用示例

### 获取客户端并执行操作

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/taouniverse/tao-mongodb"
)

func main() {
    // 获取默认实例
    client, err := mongodb.Client()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Ping 测试
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("MongoDB 连接成功")
}
```

### 数据库操作

```go
client, _ := mongodb.Client()
ctx := context.Background()

// 获取数据库和集合
db := client.Database("mydb")
coll := db.Collection("users")

// 插入文档
_, err := coll.InsertOne(ctx, map[string]interface{}{
    "name": "tao",
    "age": 18,
})

// 查询文档
var result map[string]interface{}
err = coll.FindOne(ctx, map[string]interface{}{"name": "tao"}).Decode(&result)
```

### 多实例使用

```go
// 获取主库连接（读写）
master, _ := mongodb.GetClient("master")

// 获取从库连接（只读）
replica, _ := mongodb.GetClient("replica")

// 主库写入
master.Database("mydb").Collection("logs").InsertOne(ctx, logEntry)

// 从库查询
replica.Database("mydb").Collection("logs").Find(ctx, filter)
```

## 单元测试

### 快速测试（无需 Docker）

```bash
# 仅运行配置相关测试
go test -v -run "TestConfig" ./...
```

### 完整集成测试（需要 Docker）

```bash
# 启动 MongoDB 并运行单实例测试
make test

# 启动 MongoDB 并运行多实例测试
make test-multi

# 启动 MongoDB 并运行所有测试
make test-all

# 生成覆盖率报告
make coverage

# 停止 MongoDB 服务
make down
```

### 手动测试

```bash
# 1. 启动 MongoDB
docker-compose up -d

# 2. 运行单实例测试
go test -v ./...

# 3. 运行多实例测试
TAO_TEST_MULTI_INSTANCE=true go test -v ./...

# 4. 停止 MongoDB
docker-compose down
```

## 开发指南

| 文件 | 说明 |
|------|------|
| `config.go` | InstanceConfig 字段 + ValidSelf 默认值 |
| `mongodb.go` | NewMongoDB 构造器 + 工厂注册 |
