package util

import (
	"reflect"
	"strings"
)

// InSlice haystack slice 中に needle が含まれるか
func InSlice(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// SearchDefaultCharaset search DEFAULT CHARSET={charaset}
func SearchDefaultCharaset(target string) (result string) {
	result = target[strings.LastIndex(target, "DEFAULT CHARSET="):]
	result = strings.Replace(result, "DEFAULT CHARSET=", "", -1)
	return strings.Split(result, " ")[0]
}

// TrimUnsigned trim <type> unsigned
func TrimUnsigned(target string) (result string) {
	return strings.Replace(target, " unsigned", "", -1)
}

// ContainsAutoIncrement contains auto_increment
func ContainsAutoIncrement(target string) bool {
	return strings.Contains(target, "auto_increment")
}

// MapDiffKeys Compares the keys from before map against the keys from after map and returns the difference keys
func MapDiffKeys(b, a interface{}) (rs []string) {
	bKeys := reflect.ValueOf(b).MapKeys()
	aKeys := reflect.ValueOf(a).MapKeys()

	for _, bKey := range bKeys {
		bk := bKey.String()

		isFound := false
		for _, aKey := range aKeys {
			if bk == aKey.String() {
				isFound = true
				break
			}
		}

		if !isFound {
			rs = append(rs, bk)
		}
	}

	return
}
