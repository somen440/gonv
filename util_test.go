package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInSlice(t *testing.T) {
	tests := []struct {
		title    string
		expect   bool
		needle   string
		haystack []string
	}{
		{"success", true, "a", []string{"a", "b", "c"}},
		{"failed", false, "d", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			actual := InSlice(tt.needle, tt.haystack)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestSearchDefaultCharaset(t *testing.T) {
	target := `CREATE TABLE sample (
id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
name varchar(255) NOT NULL,
created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
modified datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
`
	expected := "utf8mb4"
	actual := SearchDefaultCharaset(target)
	assert.Equal(t, expected, actual)
}
