/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */


package datasources

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

var o orm.Ormer

//初始化数据源
func InitDSs() {
	ds := config.GetProp("db", "datasources")
	if len(ds) > 0 {
		dsArr := strings.Split(ds, ",")
		for i := range dsArr {
			initDS(dsArr[i])
		}
	}
}

// 初始化数据源
func initDS(dbName string) {

	var dns string
	db_type := config.GetProp("db", "type")
	db_host := config.GetProp("db", dbName+"_host")
	db_port := config.GetProp("db", dbName+"_port")
	db_user := config.GetProp("db", dbName+"_user")
	db_pass := config.GetProp("db", dbName+"_pass")
	db_name := config.GetProp("db", dbName+"_name")
	db_maxIdle := config.GetProp("db", "maxidle")
	db_maxConn := config.GetProp("db", "maxconn")

	beego.Info("init::<db_host=", db_host, ", db_port=", db_port, ", db_name=", db_name, ">")

	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
		beego.Info("mysql info:" + dns)
		break

	default:
		beego.Critical("Database driver is not allowed:", db_type)
	}

	maxIdle, _ := strconv.Atoi(db_maxIdle)
	maxConn, _ := strconv.Atoi(db_maxConn)
	beego.Info("regist db is:" + dbName)
	orm.RegisterDataBase(dbName, db_type, dns, maxIdle, maxConn)
	orm.Debug = true

	flag.Parse()

}
