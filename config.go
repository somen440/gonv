package main

// DbDriver database の driver
type DbDriver string

// DbDriver enum
const (
	MySQL DbDriver = "mysql"
)

// AsString DbDriver to string
func (driver *DbDriver) AsString() string {
	return string(*driver)
}

// DBConfig db につなぐのに必要な設定
type DBConfig struct {
	Driver   DbDriver
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// DataSourceName sql.Open の第二引数に該当
func (conf *DBConfig) DataSourceName() string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database
}
