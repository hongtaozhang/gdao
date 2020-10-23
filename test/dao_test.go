/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package test

import (
	"gdao/dao"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"reflect"
	"testing"
)

type TestDao struct {
	Id    int64
	Name  string
	Phone string
	Age   int
}

// set the table name
func (m *TestDao) TableName() string {
	return "test_daos"
}

func init() {
	orm.RegisterModel(new(TestDao))
}

func NewTestDao(id int64, name string, phone string, age int) TestDao {

	testDao := TestDao{
		Id:    id,
		Name:  name,
		Phone: phone,
		Age:   age,
	}
	return testDao
}

func Test_DB_Insert(t *testing.T) {

	// models.Connect()

	var dbDao = dao.NewDBDao()
	// var cacheDao = dao.NewCacheDao()

	td := NewTestDao(0, "hongtaozhang", "1821223444", 23)

	//================= read  ==============//
	id, err := dbDao.Insert(&td)
	beego.Debug("DB insert sucess, id:", id, " err:", err)

	//================= read  ==============//
	var td2 TestDao
	err = dbDao.Read(&td2, id)
	beego.Debug("DB Read sucess, obj:", td2, " err:", err)

	//================= Update  ==============//
	var str = "testdao"
	td2.Name = str
	id, err = dbDao.Update(&td2)
	beego.Debug("DB Update sucess, obj:", td2, " err:", err)

	//================= Read  ==============//

	err = dbDao.Read(&td2, id)
	beego.Debug("DB Read sucess, obj:", td2, " err:", err)
	if td2.Name == str {
		beego.Debug("DB Update sucess, obj:", td2, " err:", err)
	}

	// id, err = dbDao.Delete(&td2)
	// beego.Debug("DB Delete sucess, obj:", td2, " err:", err)

}

func Test_DB_Id_List(t *testing.T) {

	// models.Connect()

	var dbDao = dao.NewDBDao()
	// var cacheDao = dao.NewCacheDao()

	//================= GetIdList  ==============//
	var td TestDao
	var tds []TestDao

	mp := dao.NewSortMap()
	mp.Set("Name", "testdao")

	//================= GetIdList  ==============//
	ids, err := dbDao.GetIdList(&td, &tds, "Test_DB_Id_List", mp, nil)
	beego.Debug("DB GetIdList sucess, id:", ids, " err:", err)

	ids, err = dbDao.GetIdListPage(&td, &tds, "Test_DB_Id_List", mp, nil, 1, 10)
	beego.Debug("DB GetIdListPage sucess, id:", ids, " err:", err)

	//================= GetIdList  ==============//
	res, err1 := dbDao.GetList(&td, &tds, "Test_DB_Id_List", mp, nil)

	// var temp TestDao
	for i, v := range res {

		var temp3 *TestDao
		temp3 = GetTestDao(v)
		beego.Debug("DB GetIdList sucess, res i:", i, " value:", v, " err:", err1, " value:", temp3.Age)
	}
	beego.Debug("DB GetIdList sucess, id:", res, " err:", err1, " type:", reflect.ValueOf(res))

	//================= GetMapingId  ==============//
	id, err3 := dbDao.GetMapingId(&td, "GetMapingId", mp)
	beego.Debug("DB GetMapingId sucess, id:", id, " err:", err3)

	//================= GetMaping  ==============//
	res4, err4 := dbDao.GetMaping(&td, "GetMapingId", mp)
	var temp3 *TestDao
	temp3 = GetTestDao(res4)
	beego.Debug("DB GetMapingId sucess, id:", temp3.Age, " err:", err4)

}

func Test_Cache_Id_List(t *testing.T) {

	// models.Connect()

	var cacheDao = dao.NewCacheDao()

	//================= GetIdList  ==============//
	var td TestDao
	var tds []TestDao

	mp := dao.NewSortMap()
	mp.Set("Name", "testdao")

	//================= GetIdList  ==============//
	ids, err := cacheDao.GetIdList(&td, &tds, "Test_Cache_Id_List", mp, nil)
	beego.Debug("Cache GetIdList sucess, id:", ids, " err:", err)

	ids, err = cacheDao.GetIdListPage(&td, &tds, "Test_Cache_Id_List", mp, nil, 1, 10)
	beego.Debug("Cache GetIdListPage sucess, id:", ids, " err:", err)

	//================= GetIdList  ==============//
	res, err1 := cacheDao.GetList(&td, &tds, "Test_Cache_Id_List", mp, nil)

	// var temp TestDao
	for i, v := range res {

		var temp3 *TestDao
		temp3 = GetTestDao(v)
		beego.Debug("Cache GetIdList sucess, res i:", i, " value:", v, " err:", err1, " value:", temp3.Age)
	}
	beego.Debug("Cache GetIdList sucess, id:", res, " err:", err1, " type:", reflect.ValueOf(res))

	//================= GetMapingId  ==============//
	id, err3 := cacheDao.GetMapingId(&td, "GetMapingId", mp)
	beego.Debug("Cache GetMapingId sucess, id:", id, " err:", err3)

	//================= GetMaping  ==============//
	res4, err4 := cacheDao.GetMaping(&td, "GetMapingId", mp)
	var temp3 *TestDao
	temp3 = GetTestDao(res4)
	beego.Debug("Cache GetMapingId sucess, id:", temp3.Age, " err:", err4)

}

func Test_MGet_Cache_Id_List(t *testing.T) {

	// models.Connect()

	var cacheDao = dao.NewCacheDao()

	//================= GetIdList  ==============//
	var td TestDao
	var tds []TestDao

	mp := dao.NewSortMap()
	mp.Set("Name", "testdao")

	//================= GetIdList  ==============//
	res, err1 := cacheDao.GetList(&td, &tds, "Test_Cache_Id_List", mp, nil)

	// var temp TestDao
	for i, v := range res {

		var temp3 *TestDao
		temp3 = GetTestDao(v)
		beego.Debug("Cache GetIdList sucess, res i:", i, " value:", v, " err:", err1, " value:", temp3.Age)
	}
}

func Test_RedisCache_Id_List(t *testing.T) {

	// utils.NewRedisClusterCache()

}
func GetTestDao(v interface{}) *TestDao {

	switch result := v.(type) {
	case *TestDao:
		beego.Debug("  assert:", result)
		return (*TestDao)(result)
	default:
		beego.Debug("  assert:", result, " type:", reflect.TypeOf(result))

	}

	return &TestDao{}
}
