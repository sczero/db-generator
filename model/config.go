package model

type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Protocol      string `json:"protocol"`
	Address       string `json:"address"`
	Dbname        string `json:"dbname"`
	TableName     string `json:"tableName"`
	OutputDir     string `json:"outputDir"`
	OutputPackage string `json:"outputPackage"`
}
