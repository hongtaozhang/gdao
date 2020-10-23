/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package dao

import (
	"fmt"
	"reflect"
	"strings"
)

func getMdReflectInfo(md interface{}, needPtr bool) (ind reflect.Value, typ reflect.Type) {

	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ = ind.Type()

	if needPtr && val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<CacheDao> cannot use non-ptr model struct"))
	}
	return ind, typ
}

//
func getMdFullName(ind reflect.Value) string {
	typ := ind.Type()
	name := typ.Name()
	return name

}

func getIdValue(ind reflect.Value) int64 {
	idValue := ind.FieldByName("Id")
	return idValue.Int()
}

func setIdValue(ind reflect.Value, id int64) {

	idValue := ind.FieldByName("Id")
	idValue.SetInt(id)

}

func cacheKey(modelName string, id int64) string {
	key := "mks_" + modelName + "_" + toString(id)
	return key
}

func cacheListKey(name string, values []interface{}) string {

	key := "listks_" + name
	keyParams := make([]string, 0)
	for _, v := range values {
		keyParams = append(keyParams, toString(v))
	}
	jsonstr := strings.Join(keyParams, "_")
	res := key + "_" + jsonstr
	return res
}

// 分页使用
func subInt64Slice(sl []int64, from int, to int) []int64 {

	if sl == nil {
		return nil
	}

	le := len(sl)
	if le == 0 {
		return nil
	}

	if from > le {
		return nil
	}

	if from < 0 || to < 0 {
		from = 0
	}

	if to > le {
		return sl[from:le]
	}
	return sl[from:to]

}

func toString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}

func offsetNano(beginNano int64, endNano int64) int64 {

	offset := endNano - beginNano
	return offset
}

func offsetMillisecond(beginNano int64, endNano int64) float64 {

	offset := float64(endNano - beginNano)
	ms := (offset / 1e6)
	return ms
}
