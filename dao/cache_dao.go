/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package dao

import (
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
	"sync"
)

type CacheDao struct {
	cache cache.Cache
	dbDao *DBDao
	sync.Mutex
}

var MONTH_SECOND_TIME = int64(30 * 24 * 60 * 60)

var cacheDao *CacheDao

var cacheDaoMap = NewSortMap()

// 采取单例模式
func NewCacheDaoByParam(datasource string) *CacheDao {

	cacheDao := cacheDaoMap.Get(datasource)
	if nil != cacheDao {
		return GetCacheDao(cacheDao)
	}

	newCacheDao := &CacheDao{
		dbDao: NewDBDaoByParam(datasource),
	}
	// 每次获取，保证在redis切换ip后能获取到最新可用redis
	newCacheDao.cache = ds.GetRedisClusterCache()
	cacheDaoMap.Set(datasource, newCacheDao)
	return newCacheDao
}

// 采取单例模式
func NewCacheDao() *CacheDao {

	cacheDao := cacheDaoMap.Get(DEFUALT)
	if nil != cacheDao {
		return GetCacheDao(cacheDao)
	}

	newCacheDao := &CacheDao{
		dbDao: NewDBDao(),
	}

	// 每次获取，保证在redis切换ip后能获取到最新可用redis
	newCacheDao.cache = ds.GetRedisClusterCache()
	cacheDaoMap.Set(DEFUALT, newCacheDao)

	return newCacheDao
}

func GetCacheDao(v interface{}) *CacheDao {

	switch result := v.(type) {
	case *CacheDao:

		return (*CacheDao)(result)
	default:
		return nil

	}
	return nil
}

func (cacheDao *CacheDao) Insert(obj interface{}) (int64, error) {

	id, err := cacheDao.dbDao.Insert(obj)

	// 获取对象反射信息
	ind, _ := getMdReflectInfo(obj, true)

	// put to cache
	if id > 0 {

		// 利用反射设置id的值,防止数据库自增id没有值.
		setIdValue(ind, id)
		mdName := getMdFullName(ind)
		jsonstr, err := utils.GetJsonFromV(obj)
		if err != nil {
			fmt.Errorf("<CacheDao> Insert convert struct to json data error, error is:", err)
		} else {
			// cacheDao.cache.Set(cacheKey(mdName, id), jsonstr)
			cacheDao.cache.Put(cacheKey(mdName, id), jsonstr, MONTH_SECOND_TIME)
		}

	}
	return id, err
}

// 先读cache
// 如果cache没有数据，从db加载，然后set到cache中
func (cacheDao *CacheDao) Read(obj interface{}, id int64) error {

	ind, _ := getMdReflectInfo(obj, true)
	mdName := getMdFullName(ind)

	// 从cache中获取数据
	cacheValue, _ := cacheDao.cache.Get(cacheKey(mdName, id))

	if cacheValue != nil {
		err := utils.GetVFromJson(utils.GetByteArray(cacheValue), obj)
		if err != nil {
			fmt.Errorf("<CacheDao> Read convert struct to json data error, error is:", err)
		}
	} else {

		//Cache 中数据不存在，则从db中读取
		err := cacheDao.dbDao.Read(obj, id)
		if err != nil {
			return err
		} else if nil != obj {
			// 写入cache
			jsonstr, err := utils.GetJsonFromV(obj)
			if err != nil {
				fmt.Errorf("<CacheDao> Read convert struct to json data error, error is:", err)
			} else {
				// cacheDao.cache.Set(cacheKey(mdName, id), jsonstr)
				cacheDao.cache.Put(cacheKey(mdName, id), jsonstr, MONTH_SECOND_TIME)
			}

		}
	}
	return nil
}

// 先读cache
// 如果cache没有数据，从db加载，然后set到cache中
func (cacheDao *CacheDao) MRead(obj interface{}, ids []int64) ([]interface{}, error) {

	ind, _ := getMdReflectInfo(obj, true)
	mdName := getMdFullName(ind)

	keys := make([]string, len(ids))
	keysMap := make(map[int64]int64, len(ids))
	for i, v := range ids {
		key := cacheKey(mdName, v)
		keysMap[v] = v
		keys[i] = key

	}

	objs := make(map[int64]interface{}, 0)
	// 从cache中获取数据
	cacheValues, _ := cacheDao.cache.MGet(keys)
	if cacheValues != nil {
		for _, v := range cacheValues {

			obb := reflect.New(ind.Type()).Interface()
			err := utils.GetVFromJson(utils.GetByteArray(v), obb)
			if err != nil {
				continue
			}

			// 已经在cache中存在删除
			tind, _ := getMdReflectInfo(obb, true)
			id := getIdValue(tind)
			delete(keysMap, id)
			objs[id] = obb
		}
	}

	// 查询cache中不存在的id
	if len(keysMap) > 0 {

		for _, v := range keysMap {
			obb := reflect.New(ind.Type()).Interface()
			//Cache 中数据不存在，则从db中读取
			err := cacheDao.dbDao.Read(obb, v)
			if err != nil {
				beego.Error("Read Data error, id", v, " err info:", err)
				continue
			} else {
				// 写入cache
				jsonstr, err := utils.GetJsonFromV(obb)

				if err != nil {
					fmt.Errorf("<CacheDao> Read convert struct to json data error, error is:", err)
				} else if nil != obb {
					// cacheDao.cache.Set(cacheKey(mdName, v), jsonstr)
					cacheDao.cache.Put(cacheKey(mdName, v), jsonstr, MONTH_SECOND_TIME)
					tind, _ := getMdReflectInfo(obb, true)
					id := getIdValue(tind)
					delete(keysMap, id)
					objs[id] = obb
				}
			}
		}
	}

	// 保证顺序，重新遍历ids列表
	obres := make([]interface{}, 0)
	for _, tid := range ids {
		t := objs[tid]
		if t == nil {
			continue
		}
		obres = append(obres, t)
	}

	return obres, nil
}

func (cacheDao *CacheDao) Update(obj interface{}) (int64, error) {

	num, err := cacheDao.dbDao.Update(obj)

	//fmt.Println("====================num:", num, "  obj:", obj)
	// 获取对象反射信息
	ind, _ := getMdReflectInfo(obj, true)

	// 数据放入缓存
	if num > 0 {
		mdName := getMdFullName(ind)
		id := getIdValue(ind)
		jsonstr, err := utils.GetJsonFromV(obj)
		if err != nil {
			fmt.Errorf("<CacheDao> Update convert struct to json data error, error is:", err)
		} else if nil != obj {
			// cacheDao.cache.Set(cacheKey(mdName, id), jsonstr)
			cacheDao.cache.Put(cacheKey(mdName, id), jsonstr, MONTH_SECOND_TIME)
		}

	}
	return num, err
}

func (cacheDao *CacheDao) Delete(obj interface{}) (int64, error) {

	// 先删除cache
	ind, _ := getMdReflectInfo(obj, true)
	mdName := getMdFullName(ind)
	id := getIdValue(ind)
	cacheDao.cache.Delete(cacheKey(mdName, id))

	// 数据库中删除
	id, err := cacheDao.dbDao.Delete(obj)
	return id, err
}

func (cacheDao *CacheDao) Count(obj interface{}, name string, params *SortMap) (int64, error) {
	return cacheDao.dbDao.Count(obj, name, params)
}

//查询cache中是否存在某个key
func (cacheDao *CacheDao) IsExistInCache(key string) (bool, error) {

	isOk, err := cacheDao.cache.IsExist(key)
	return isOk, err
}

func (cacheDao *CacheDao) DelKeyFromCache(key string) (bool, error) {
	err := cacheDao.cache.Delete(key)
	if nil != err {
		return false, err
	}
	return true, nil
}

func (cacheDao *CacheDao) GetIdList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]int64, error) {

	// 查询缓存
	cacheValue, _ := cacheDao.cache.Get(cacheListKey(name, params.ValueSet()))
	var ids []int64
	if cacheValue != nil {
		err := utils.GetVFromJson(utils.GetByteArray(cacheValue), &ids)
		if err != nil {
			fmt.Errorf("<CacheDao> GetIdList convert struct to json data error, error is:", err)
		}
		return ids, nil
	} else {

		// put to cache
		//Cache 中数据不存在，则从db中读取
		ids, err := cacheDao.dbDao.GetIdList(obj, res, name, params, orderby)
		if err != nil {
			return nil, err
		} else if nil != ids && len(ids) > 0 {
			// 写入cache
			jsonstr, err := utils.GetJsonFromV(ids)
			if err != nil {
				fmt.Errorf("<CacheDao> GetIdList convert struct to json data error, error is:", err)
			} else {
				// cacheDao.cache.Set(cacheListKey(name, params.ValueSet()), jsonstr)
				cacheDao.cache.Put(cacheListKey(name, params.ValueSet()), jsonstr, MONTH_SECOND_TIME)
			}
		}
		return ids, nil
	}
	return nil, nil
}

func (cacheDao *CacheDao) GetIdListPage(obj interface{}, res interface{}, name string, params *SortMap, orderby []string, start int64, count int64) ([]int64, error) {

	// 查询缓存
	cacheValue, _ := cacheDao.cache.Get(cacheListKey(name, params.ValueSet()))
	var ids []int64

	if cacheValue != nil {
		err := utils.GetVFromJson(utils.GetByteArray(cacheValue), &ids)
		if err != nil {
			fmt.Errorf("<CacheDao> GetIdListPage convert struct to json data error, error is:", err)
		}
		// 分页获取ids
		if count < 0 {
			return ids, nil
		}
		subIds := subInt64Slice(ids, int(start), int(start+count))
		return subIds, nil

	} else {

		// put to cache
		//Cache 中数据不存在，则从db中读取
		tcount, errc := cacheDao.dbDao.Count(obj, name, params)
		if errc != nil {
			fmt.Errorf("<CacheDao> GetIdListPage query db count error, error is:", errc)
		}
		ids, err := cacheDao.dbDao.GetIdListPage(obj, res, name, params, orderby, 0, tcount)
		if err != nil {
			return nil, err
		} else if nil != ids && len(ids) > 0 {
			// 写入cache
			jsonstr, err := utils.GetJsonFromV(ids)
			if err != nil {
				fmt.Errorf("<CacheDao> GetIdListPage convert struct to json data error, error is:", err)
			} else {
				// cacheDao.cache.Set(cacheListKey(name, params.ValueSet()), jsonstr)
				cacheDao.cache.Put(cacheListKey(name, params.ValueSet()), jsonstr, MONTH_SECOND_TIME)
			}
		}
		// 分页获取ids
		if count < 0 {
			return ids, nil
		}
		subIds := subInt64Slice(ids, int(start), int(start+count))
		return subIds, err
	}
}

func (cacheDao *CacheDao) GetListByIds(obj interface{}, ids []int64) ([]interface{}, error) {

	// 没有数据情况
	if len(ids) <= 0 {
		err := fmt.Errorf("<CacheDao> GetListByIds param ids can not be empty!")
		return nil, err
	}

	// var resSlice = make([]interface{}, 0)
	// val := reflect.ValueOf(obj)
	// ind := reflect.Indirect(val)

	// for _, id := range ids {

	// 	obb := reflect.New(ind.Type()).Interface()
	// 	cacheDao.Read(obb, id)
	// 	resSlice = append(resSlice, obb)
	// }

	// return resSlice, nil

	value, err := cacheDao.MRead(obj, ids)
	return value, err
}

func (cacheDao *CacheDao) GetList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]interface{}, error) {

	ids, err := cacheDao.GetIdList(obj, res, name, params, orderby)
	if err != nil {
		return nil, err
	}

	return cacheDao.GetListByIds(obj, ids)
}

// 分页根据查询条件获取对象列表
// 返回具体对象
func (cacheDao *CacheDao) GetListPage(obj interface{}, res interface{}, name string, params *SortMap, orderby []string, start int64, count int64) ([]interface{}, error) {

	ids, err := cacheDao.GetIdListPage(obj, res, name, params, orderby, start, count)
	if err != nil {
		return nil, err
	}
	return cacheDao.GetListByIds(obj, ids)
}

func (cacheDao *CacheDao) GetMapingId(obj interface{}, name string, params *SortMap) (int64, error) {

	// 查询缓存
	cacheValue, _ := cacheDao.cache.Get(cacheListKey(name, params.ValueSet()))

	var id int64
	if cacheValue != nil {
		err := utils.GetVFromJson(utils.GetByteArray(cacheValue), &id)
		if err != nil {
			fmt.Errorf("<CacheDao> GetMapingId convert struct to json data error, error is:", err)
		}
		return id, nil
	} else {

		// put to cache
		//Cache 中数据不存在，则从db中读取
		id, err := cacheDao.dbDao.GetMapingId(obj, name, params)
		if err != nil {
			return 0, err
		} else {
			// 写入cache
			jsonstr, err := utils.GetJsonFromV(id)
			if err != nil {
				fmt.Errorf("<CacheDao> GetIdList convert struct to json data error, error is:", err)
			} else {
				// cacheDao.cache.Set(cacheListKey(name, params.ValueSet()), jsonstr)
				cacheDao.cache.Put(cacheListKey(name, params.ValueSet()), jsonstr, MONTH_SECOND_TIME)

			}
		}
		return id, nil
	}
	return 0, nil
}

// 此方法适用于多个条件对应于数据库唯一一条记录场景，即map映射
// 返回唯一对象s
func (cacheDao *CacheDao) GetMaping(obj interface{}, name string, params *SortMap) (interface{}, error) {

	id, err := cacheDao.GetMapingId(obj, name, params)
	if err != nil {
		return nil, err
	}

	cacheDao.Read(obj, id)
	return obj, nil
}
