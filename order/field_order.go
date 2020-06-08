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
