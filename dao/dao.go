/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package dao

import ()

type Dao interface {

	// 插入对象
	// 返回对象Id，Id在数据为此对象对应的表的主键
	Insert(obj interface{}) (int64, error)

	// 更新对象
	// 返回对象Id
	Update(obj interface{}) (int64, error)

	// 查询对象，注意：查询前需要设置对象Id字段值
	Read(obj interface{}, id int64) error

	// 删除对象，注意：查询前需要设置对象Id字段值
	// 返回删除对象Id
	Delete(obj interface{}) (int64, error)

	//计算记录数量
	Count(obj interface{}, name string, params *SortMap) (int64, error)

	//查询cache中是否存在某个key
	IsExistInCache(key string) (bool, error)

	//从cache中删除key
	DelKeyFromCache(key string) (bool, error)

	// 获取查询条件对应的对象Id列表
	// 返回满足查询条件的Id列表
	GetIdList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]int64, error)

	// 分页获取查询条件对应的对象Id列表
	// 分页返回满足查询条件的Id列表
	GetIdListPage(obj interface{}, res interface{}, name string, params *SortMap, orderby []string, start int64, count int64) ([]int64, error)

	// 根据查询条件获取对象列表
	// 返回具体对象
	GetList(obj interface{}, res interface{}, name string, params *SortMap, orderby []string) ([]interface{}, error)

	// 根据现有的ids列表，查询对象
	// 返回为对象列表
	GetListByIds(obj interface{}, ids int64) ([]interface{}, error)

	// 分页根据查询条件获取对象列表
	// 返回具体对象
	GetListPage(name string, params *SortMap, orderby []string, start int64, count int64) ([]interface{}, error)

	// 此方法适用于多个条件对应于数据库唯一一条记录场景，即map映射
	// 返回唯一对象的Id
	GetMapingId(obj interface{}, name string, params *SortMap) (int64, error)

	// 此方法适用于多个条件对应于数据库唯一一条记录场景，即map映射
	// 返回唯一对象
	GetMaping(obj interface{}, name string, params *SortMap) (interface{}, error)
}
