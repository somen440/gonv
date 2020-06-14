/*
Copyright 2020 somen440

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

// DataSourceName sql.Open の第二引数に該当 DSN
func (conf *DBConfig) DataSourceName() string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?charset=utf8mb4&parseTime=true"
}

// DataSourceNameNoDatabase sql.Open の第二引数に該当 DSN
func (conf *DBConfig) DataSourceNameNoDatabase() string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/"
}
