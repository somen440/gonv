package gonv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSourceName(t *testing.T) {
	expected := "user:pass@tcp(host:port)/db"

	config := &DBConfig{
		Driver:   MySQL,
		User:     "user",
		Password: "pass",
		Host:     "host",
		Port:     "port",
		Database: "db",
	}
	actual := config.DataSourceName()

	assert.Equal(t, expected, actual)
}
