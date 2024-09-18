## 快速开始

### mysql

1. 用docker启动一个mysql容器

```shell
docker run --name ca_mysql -e MYSQL_ROOT_PASSWORD=root -d -e MYSQL_DATABASE=ca_mysql -p 8086:3306 mysql:8.0
```

2. 进入mysql容器

```shell
docker ps
docker exec -it <CONTAINER ID> mysql -uroot -proot
```

3. 创建表

```sql
use
ca_mysql;
CREATE TABLE IF NOT EXISTS t_user
(
    `id`
    INT
    NOT
    NULL
    AUTO_INCREMENT,
    `name`
    VARCHAR
(
    100
) NOT NULL COMMENT '姓名',
    `password` varchar
(
    255
) NOT NULL DEFAULT '' COMMENT '密码',
    `nickname` varchar
(
    100
) NOT NULL DEFAULT '' COMMENT '昵称',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `creator` VARCHAR
(
    100
) NOT NULL DEFAULT '' COMMENT '创建人',
    `modify_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次修改时间',
    `modifier` VARCHAR
(
    100
) NOT NULL DEFAULT '' COMMENT '最后一次修改人',
    `cert` VARCHAR (5000) NULL DEFAULT '' COMMENT '证书',
    PRIMARY KEY
(
    id
)
    );
```

### redis

1. 用docker启动一个redis容器

```shell
docker run --name ca_redis -d -p 8089:6379 redis:7.4.0
```

2. 进入redis容器

```shell
docker ps
docker exec -it <CONTAINER ID> redis-cli -h 0.0.0.0 -p 6379
```

### 启动项目

```shell
cd cmd
go run main.go
```

用户注册：`http://localhost:8080/static/register.html`
用户登录：`http://localhost:8080/static/login.html`
