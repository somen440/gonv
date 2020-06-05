package gonv

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

// DataSourceName sql.Open の第二引数に該当 DSN
func (conf *DBConfig) DataSourceName() string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database
}

// DataSourceNameNoDatabase sql.Open の第二引数に該当 DSN
func (conf *DBConfig) DataSourceNameNoDatabase() string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/"
}
