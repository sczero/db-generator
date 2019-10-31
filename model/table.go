package model

type Table struct {
	TableName      string `db:"table_name"`
	TableComment   string `db:"table_comment"`
}
