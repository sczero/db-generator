# go-db-generator

## 用途说明

生成 go 的 struct(只支持mysql)

## 使用方法
```shell script
db-generator config.template.json
```

```json
{
  "username": "用户名",
  "password": "密码",
  "protocol": "tcp",
  "address": "127.0.0.1:3306",
  "dbname": "数据库名称",
  "tableName": "表名(可空)",
  "outputDir": "输出目录",
  "outputPackage": "struct文件的包名"
}
```