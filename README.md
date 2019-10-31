# go-db-generator

## 作用

生成 go 的 struct,现在只支持 mysql

## 使用方法


```go
//run.go配置
var useruame = ""
var password = ""
var protocol = "tcp"
var address = "127.0.0.1:3360"
var dbname = ""
var tableName = "" //可空,默认生成数据库下所有的表
var outputDir = "" //输出的路径
var outputPackage = "" //文件的包名
```