# 本项目用于Gin后端学习
---
## 1.运行环境
---
|组件|版本|
|:---:|:---:|
|Go|1.27.4|
|Gin|v1.12.0|
|Gorm|v1.31.2|
|Mysql|8.0.46|
|JWT| `github.com/golang-jwt/jwt` v3.2.2|
|bcrypt| `golang.org/x/crypto/bcrypt` |
|viper|spf13/viper(YAML)|

## 2.快速开始
---
### 2.1 准备数据库
```sql
CREATE DATABASE go_mall_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

### 2.2 准备配置
> config/config.example.yml提供了配置模板
```bash
cp config/config.example.yml config/config.yml
```

编辑config/config.yaml, 填入实际mysql信息:
```yaml
app:
  name: GoMallApp
  port: :3000

database:
  host: 127.0.0.1
  port: 3306
  user: <db-user>
  password: <db-password>
  name: go_mall_db
  maxidleconns: 10
  maxopenconns: 100
```

### 2.3 启动
go mod download
go run .

启动成功后访问<http://localhost:3000/>以访问Web页面
>首次启动会通过 GORM `AutoMigrate` 自动在数据库中创建 `user` 和 `todo_lists` 两张表

### 2.4 停止
使用Ctrl+C以终止整个进
