package model

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func Init() { //大写开头表示公共方法 属于导出函数 能够被其他包引用 执行数据库全局初始化
	// 获取数据库配置
	dbHost, _ := beego.AppConfig.String("db_host")
	dbPort, _ := beego.AppConfig.String("db_port")
	dbUser, _ := beego.AppConfig.String("db_user")
	dbPassword, _ := beego.AppConfig.String("db_password")
	dbName, _ := beego.AppConfig.String("db_name")
	dbCharset, _ := beego.AppConfig.String("db_charset")

	// 构建数据库连接字符串
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset)

	// 注册数据库
	err := orm.RegisterDataBase("default", "mysql", dataSource) //default是数据库别名,在允许连接多个数据库的情形下表示系统预设的默认别名
	if err != nil {
		panic(err)
	}

	// 注册所有模型
	orm.RegisterModel(
		new(User),
		new(UserPreferences),
		new(Tag),
		new(UserTag),
	)

	// 开发模式下自动创建表
	if beego.BConfig.RunMode == "dev" {
		err = orm.RunSyncdb("default", false, true)
		if err != nil {
			panic(err)
		}
	}
}
