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

package util

import (
	"reflect"
	"testing"

	"github.com/somen440/gonv/structure"
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

func TestTrimUnsigned(t *testing.T) {
	tests := []struct {
		target   string
		expected string
	}{
		{"bigint(20) unsigned", "bigint(20)"},
		{"varchar(255)", "varchar(255)"},
		{"datetime", "datetime"},
	}
	for _, tt := range tests {
		actual := TrimUnsigned(tt.target)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestContainsAutoIncrement(t *testing.T) {
	tests := []struct {
		target   string
		expected bool
	}{
		{"auto_increment", true},
		{"", false},
	}
	for _, tt := range tests {
		actual := ContainsAutoIncrement(tt.target)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestDeepEqualStruct(t *testing.T) {
	type User struct {
		Name string
	}
	a := &User{
		Name: "a",
	}
	b := User{
		Name: "a",
	}
	assert.True(t, reflect.DeepEqual(a, &b))
}

func TestMapDiffKey(t *testing.T) {
	bm := map[structure.IndexKey]*structure.IndexStructure{
		structure.IndexKey("key1"): structure.NewIndexStructure(
			structure.IndexKey("key1"),
			"type",
			true,
			[]string{"key1"},
			0,
		),
		structure.IndexKey("key2"): structure.NewIndexStructure(
			structure.IndexKey("key2"),
			"type",
			true,
			[]string{"key2"},
			1,
		),
	}
	am := map[structure.IndexKey]*structure.IndexStructure{
		structure.IndexKey("key1"): structure.NewIndexStructure(
			structure.IndexKey("key1"),
			"type",
			true,
			[]string{"key1"},
			0,
		),
	}

	expected := []string{"key2"}
	actual := MapDiffKeys(bm, am)

	assert.Equal(t, expected, actual)
}
