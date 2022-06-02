package models

type DataSourceType string

const (
	Sql         DataSourceType = "Sql"
	HttpRequest DataSourceType = "HttpRequest"
	HttpServer  DataSourceType = "HttpServer"
	RandomValue DataSourceType = "RandomValue"
	Callback    DataSourceType = "Callback"
)
