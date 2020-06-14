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

package order

import "reflect"

// FieldOrder order
type FieldOrder struct {
	Field              string
	PreviousAfterField string
	NextAfterField     string
}

// FieldOrderMap field order map
type FieldOrderMap map[string]*FieldOrder

// GenerateFieldOrderList hoge
func GenerateFieldOrderList(beforeList, afterList []string) FieldOrderMap {
	results := FieldOrderMap{}

	if reflect.DeepEqual(beforeList, afterList) {
		return results
	}

	toOrder := func(list []string) []*FieldOrder {
		results := []*FieldOrder{}
		for i, v := range list {
			var previous, next string
			if i == 0 {
				previous = ""
			} else {
				previous = list[i-1]

			}
			if len(list) == i+1 {
				next = ""
			} else {
				next = list[i+1]
			}
			results = append(results, &FieldOrder{
				Field:              v,
				PreviousAfterField: previous,
				NextAfterField:     next,
			})
		}
		return results
	}
	beforeOrderList := toOrder(beforeList)
	afterOrderList := toOrder(afterList)

	for i, beforeOrder := range beforeOrderList {
		afterOrder := afterOrderList[i]

		if beforeOrder.Field == afterOrder.Field {
			continue
		}

		find := func(targets []*FieldOrder, search string) *FieldOrder {
			for _, target := range targets {
				if search == target.Field {
					return target
				}
			}
			return nil
		}

		findOrder := find(afterOrderList, beforeOrder.Field)
		if !(beforeOrder.NextAfterField != findOrder.NextAfterField &&
			beforeOrder.PreviousAfterField != findOrder.PreviousAfterField) {
			continue
		}

		move := beforeOrder.Field
		previous := beforeOrder.PreviousAfterField

		var next string
		for _, afterOrder := range afterOrderList {
			if move == afterOrder.Field {
				next = afterOrder.PreviousAfterField
				break
			}
		}

		results[move] = &FieldOrder{
			Field:              move,
			PreviousAfterField: previous,
			NextAfterField:     next,
		}
	}

	return results
}
