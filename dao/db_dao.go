/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package dao

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
)

type DBDao struct {
	o orm.Ormer
}

const (
	DEFUALT = "default"
	PHEONIX = "pheonix"
)

var daoMap = NewSortMap()

// 采取单例模式
func NewDBDao() *DBDao {

	dbDao := daoMap.Get(DEFUALT)
	if nil != dbDao {
		return GetDBDao(dbDao)
	}

	newDbDao := &DBDao{
		o: ds.GetOrm(),
	}
	daoMap.Set(DEFUALT, newDbDao)
	return newDbDao
}

// 采取单例模式
func NewDBDaoByParam(datasource string) *DBDao {

	dbDao := daoMap.Get(datasource)
	if nil != dbDao {
		return GetDBDao(dbDao)
	}

	var newDbDao *DBDao
	if datasource == DEFUALT {
		newDbDao = &DBDao{
			o: ds.GetOrm(),
		}
		daoMap.Set(DEFUALT, newDbDao)
	} else if datasource == PHEONIX {
		newDbDao = &DBDao{
			o: ds.GetPhoenixOrm(),
		}
		daoMap.Set(PHEONIX, newDbDao)
	}
	return newDbDao
}

func GetDBDao(v interface{}) *DBDao {

	switch result := v.(type) {
	case *DBDao:

		return (*DBDao)(result)
	default:
		return nil

	}
	return nil
}

func (dao *DBDao) Insert(obj interface{}) (int64, error) {
	id, err := dao.o.Insert(obj)
	return id, err
}

func (dao *DBDao) Update(obj interface{}) (int64, error) {
	id, err := dao.o.Update(obj)
	return id, err
}

func (dao *DBDao) Read(obj interface{}, id int64) error {

	err := dao.o.QueryTable(obj).Filter("Id", id).One(obj)
	return err
}

func (dao *DBDao) Delete(obj interface{}) (int64, error) {
	id, err := dao.o.Delete(obj)
	return id, err
}

func (dao *DBDao) Count(obj interface{}, name string, params *SortMap) (int64, error) {

	// 判断查询条件
	keys := params.KeySet()
	qs := dao.o.QueryTable(obj)

	// 添加过滤条件
	for _, k := range keys {
		value := params.Values[k]
		qs = qs.Filter(toString(k), value)
	}
	count, err := qs.Count()
	return count, err
}

//查询cache中是否存在某个key
func (dao *DBDao) IsExistInCache(key string) (bool, error) {
	return false, nil
}

func (dao *DBDao) DelKeyFromCache(key string) (bool, error) {
	return false, nil
}

func (dao *DBDao) GetIdList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]int64, error) {

	// 判断查询条件
	keys := params.KeySet()
	qs := dao.o.QueryTable(obj)

	// 添加过滤条件
	for _, k := range keys {
		value := params.Values[k]
		qs = qs.Filter(toString(k), value)
	}

	// 添加排序过滤
	for _, v := range orderby {
		qs = qs.OrderBy(v)
	}

	qs.All(res, "Id")

	val := reflect.ValueOf(res)
	ind := reflect.Indirect(val)

	var ids = make([]int64, 0)
	for i := 0; i < ind.Len(); i++ {
		si := ind.Index(i)
		idv := si.FieldByName("Id").Int()
		ids = append(ids, idv)
	}

	return ids, nil
}

func (dao *DBDao) GetIdListPage(obj interface{}, res interface{}, name string, params *SortMap, orderby []string, start int64, count int64) ([]int64, error) {

	// 判断查询条件
	keys := params.KeySet()
	qs := dao.o.QueryTable(obj)

	// 添加过滤条件
	for _, k := range keys {
		value := params.Values[k]
		qs = qs.Filter(toString(k), value)
	}

	// 添加排序过滤
	for _, v := range orderby {
		qs = qs.OrderBy(v)
	}

	qs.Limit(count, start).All(res, "Id")

	val := reflect.ValueOf(res)
	ind := reflect.Indirect(val)

	var ids = make([]int64, 0)
	for i := 0; i < ind.Len(); i++ {
		si := ind.Index(i)
		idv := si.FieldByName("Id").Int()
		ids = append(ids, idv)
	}
	return ids, nil
}

func (dao *DBDao) GetListByIds(obj interface{}, ids []int64) ([]interface{}, error) {

	// 没有数据情况
	if len(ids) <= 0 {
		err := fmt.Errorf("Param ids can not be empty!")
		return nil, err
	}

	var resSlice = make([]interface{}, 0)
	val := reflect.ValueOf(obj)
	ind := reflect.Indirect(val)

	for _, id := range ids {

		obb := reflect.New(ind.Type()).Interface()
		dao.Read(obb, id)
		resSlice = append(resSlice, obb)
	}

	return resSlice, nil
}

func (dao *DBDao) GetList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]interface{}, error) {

	ids, err := dao.GetIdList(obj, res, name, params, orderby)
	if err != nil {
		panic(fmt.Errorf("<DBDao> GetList error, error is:", err))
	}
	return dao.GetListByIds(obj, ids)

}

// 分页根据查询条件获取对象列表
// 返回具体对象
func (dao *DBDao) GetListPage(obj interface{}, res interface{}, name string, params *SortMap, orderby []string, start int64, count int64) ([]interface{}, error) {
	ids, err := dao.GetIdListPage(obj, res, name, params, orderby, start, count)
	if err != nil {
		panic(fmt.Errorf("<DBDao> GetListPage error, error is:", err))
	}
	return dao.GetListByIds(obj, ids)

}

func (dao *DBDao) GetMapingId(obj interface{}, name string, params *SortMap) (int64, error) {

	// 判断查询条件
	keys := params.KeySet()
	qs := dao.o.QueryTable(obj)

	// 添加过滤条件
	for _, k := range keys {
		value := params.Values[k]
		qs = qs.Filter(toString(k), value)
	}

	qs.One(obj, "Id")
	val := reflect.ValueOf(obj)
	ind := reflect.Indirect(val)

	id := ind.FieldByName("Id").Int()
	return id, nil
}

// 此方法适用于多个条件对应于数据库唯一一条记录场景，即map映射
// 返回唯一对象s
func (dao *DBDao) GetMaping(obj interface{}, name string, params *SortMap) (interface{}, error) {

	id, err := dao.GetMapingId(obj, name, params)
	if err != nil {
		panic(fmt.Errorf("<DBDao> GetMaping error, error is:", err))
	}
	dao.Read(obj, id)
	return obj, nil

}
