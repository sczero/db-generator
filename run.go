package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sczero/db-generator/model"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

//DSN:username:password@protocol(address)/dbname?param=value
var useruame = ""
var password = ""
var protocol = "tcp"
var address = "127.0.0.1:3360"
var dbname = ""
var tableName = ""
var outputDir = ""
var outputPackage = ""

var pool, _ = sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", useruame, password, protocol, address, dbname))

func main() {
	e := os.MkdirAll(outputDir, os.ModePerm)
	if e != nil {
		panic(e)
	}
	tables := queryTables(tableName)

	_ = os.MkdirAll(outputDir, os.ModePerm)
	builder := strings.Builder{}
	for _, table := range tables {
		//转换表名
		builder.Reset()
		packageTime := false
		packageSql := false
		builder.WriteString(fmt.Sprintf("//%s\ntype %s struct {\n", table.TableComment, CamelStr(table.TableName)))
		//拼接字符串
		for _, column := range queryColumns(table.TableName) {
			//转换列名
			dataType := strings.ToUpper(column.DataType)
			value, ok := model.DataTypeMap[dataType]
			if ok {
				if column.IsNullable == "YES" {
					dataType = value[1]
					packageSql = true
				} else {
					dataType = value[0]
				}
				//是否需要 sql 包
				packageTime = dataType == "time.Time"

			} else {
				dataType = "string"
			}
			//拼接字符串
			camelStr := CamelStr(column.ColumnName)
			builder.WriteString(fmt.Sprintf("	%s %s `db:\"%s\";json:\"%s\"` //%s", camelStr, dataType, column.ColumnName, strings.ToLower(string(camelStr[0]))+camelStr[1:], column.ColumnComment))
			if column.ColumnKey != "" {
				builder.WriteString("(" + column.ColumnKey + ")")
			}
			builder.WriteString("\n")
		}
		builder.WriteString("}\n")
		fileStr := "package " + outputPackage + "\nimport ("
		if packageSql {
			fileStr += "\"database/sql\"\n"
		}
		if packageTime {
			fileStr += "\"time\"\n"
		}
		fileStr += ")\n"
		fileStr += builder.String()
		_ = ioutil.WriteFile(path.Join(outputDir, table.TableName+".go"), []byte(fileStr), os.ModePerm)
	}
	_ = os.Chdir(outputDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}
	fmt.Printf("格式化结果:\n%s\n", string(out))
}

//查询所有的列
func queryColumns(tableName string) []model.Column {
	var results []model.Column
	e := pool.Select(&results, "select COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT from information_schema.COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?", dbname, tableName)
	if e != nil {
		panic(e)
	}
	return results
}

//查询所有的表
func queryTables(tableName string) []model.Table {
	var tables []model.Table
	sql := "SELECT table_name ,table_comment FROM information_schema.TABLES WHERE table_schema = '" + dbname + "'"
	if tableName != "" {
		sql += " and table_name = '" + tableName + "'"
	}
	sql += " ORDER BY table_name"
	e := pool.Select(&tables, sql)
	if e != nil {
		panic(e)
	}
	return tables
}

//下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}
