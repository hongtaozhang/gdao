/**
 *
 * author: hongtaozhang
 * date  : 2020-10-23 11:29
 */

package datasources

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	InitDSs()
}

func GetOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("test")
	return o
}
