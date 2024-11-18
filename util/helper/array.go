package helper

import (
	"fmt"
	"reflect"
)

// 返回两个切片的交集
func ArrayIntersectString(a, b []string) []string {
	// 使用 map 记录元素出现情况
	exists := make(map[string]bool)
	result := []string{}
	// 遍历第一个数组，并将元素存储在 map 中
	for _, v := range a {
		exists[v] = true
	}
	// 遍历第二个数组，检查元素是否在 map 中
	for _, v := range b {
		if exists[v] {
			result = append(result, v)
		}
	}
	return result
}

// 返回两个切片的交集
func ArrayIntersectUint64(a, b []uint64) []uint64 {
	// 使用 map 记录元素出现情况
	exists := make(map[uint64]bool)
	result := []uint64{}
	// 遍历第一个数组，并将元素存储在 map 中
	for _, v := range a {
		exists[v] = true
	}
	// 遍历第二个数组，检查元素是否在 map 中
	for _, v := range b {
		if exists[v] {
			result = append(result, v)
		}
	}
	return result
}

// in_array Uint64功能
func InArrayUint64(str uint64, array []uint64) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// InArrayInt in_array int功能
func InArrayInt(str int, array []int) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// in_arrayInt64功能
func InArrayInt64(str int64, array []int64) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// in_arrayInt8功能
func InArrayInt8(str int8, array []int8) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// in_array string功能
func InArrayString(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

// ArrayColumnAny array_column功能
// ArrayColumnAny resultArr := ArrayColomn[struct_name, string](list, "Name")
// ArrayColumnAny resultArr.([]string)
func ArrayColumnAny[T, V any](array []T, k any) any {
	values := make([]V, len(array))
	switch reflect.TypeOf(array).Elem().Kind() {
	case reflect.Slice, reflect.Array:
		for i, v := range array {
			values[i] = reflect.ValueOf(v).Index(int(reflect.ValueOf(k).Int())).Interface().(V)
		}
		break
	case reflect.Map:
		for i, v := range array {
			values[i] = reflect.ValueOf(v).MapIndex(reflect.ValueOf(k)).Interface().(V)
		}
		break
	case reflect.Struct:
		for i, v := range array {
			values[i] = reflect.ValueOf(v).FieldByName(reflect.ValueOf(k).String()).Interface().(V)
		}
		break
	}
	return values
}

// ArrayColumnIndexMapStringAny array_column Index功能
// ArrayColumnIndexMapStringAny resultArr := ArrayColomn[struct_name, string](list, "Name")
// ArrayColumnIndexMapStringAny resultArr.([]string)
func ArrayColumnIndexMapStringAny[T, V any](array []T, index any, k any) any {
	var values = map[string]V{}
	switch reflect.TypeOf(array).Elem().Kind() {
	case reflect.Slice, reflect.Array:
		for _, v := range array {
			indexValue := fmt.Sprintf("%v", reflect.ValueOf(v).Index(int(reflect.ValueOf(index).Int())).Interface())
			values[indexValue] = reflect.ValueOf(v).Index(int(reflect.ValueOf(k).Int())).Interface().(V)
		}
		break
	case reflect.Map:
		for _, v := range array {
			indexValue := fmt.Sprintf("%v", reflect.ValueOf(v).MapIndex(reflect.ValueOf(index)).Interface())
			values[indexValue] = reflect.ValueOf(v).MapIndex(reflect.ValueOf(k)).Interface().(V)
		}
		break
	case reflect.Struct:
		for _, v := range array {
			indexValue := fmt.Sprintf("%v", reflect.ValueOf(v).FieldByName(reflect.ValueOf(index).String()).Interface())
			values[indexValue] = reflect.ValueOf(v).FieldByName(reflect.ValueOf(k).String()).Interface().(V)
		}
		break
	}
	return values
}

// array_merge uint64功能
func ArrayMergeUint64(ss ...[]uint64) []uint64 {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]uint64, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// array_merge string功能
func ArrayMergeString(ss ...[]string) []string {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]string, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// array_merge interface功能 resultArr := ArrayMergeAny[struct_name, string](xxx, xx)
func ArrayMergeAny[T, V any](ss ...[]T) []T {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]T, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// ArrayUnique array_unique功能
func ArrayUniqueString(array []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range array {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ArrayUnique array_unique功能
func ArrayUnique(array []interface{}) []interface{} {
	var result []interface{}
	temp := map[interface{}]struct{}{}
	for _, v := range array {
		if _, ok := temp[v]; !ok {
			temp[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// ArrayUniqueUint64 array_unique功能
func ArrayUniqueUint64(slice []uint64) []uint64 {
	keys := make(map[uint64]bool)
	var list []uint64

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
