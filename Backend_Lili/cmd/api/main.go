package main

import (
	"log"

	"Backend_Lili/internal/router"
	"Backend_Lili/internal/user/model"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

// 全局初始化
func init() {
	log.Println("=== 理理小程序后端服务启动 ===")

	// 1. 加载配置文件
	if err := loadConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 2. 初始化数据库
	if err := initDatabase(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 3. 注册路由
	if err := registerRoutes(); err != nil {
		log.Fatalf("路由注册失败: %v", err)
	}

	log.Println("=== 初始化完成 ===")
}

// loadConfig 加载配置文件
func loadConfig() error {
	log.Println("正在加载配置文件...")

	// 设置配置文件路径
	err := beego.LoadAppConfig("ini", "pkg/conf/app.conf")
	if err != nil {
		return err
	}

	// 验证关键配置
	requiredConfigs := []string{"db_host", "db_port", "db_user", "db_password", "db_name"}
	for _, config := range requiredConfigs {
		if _, err := beego.AppConfig.String(config); err != nil {
			return err
		}
	}

	log.Println("配置文件加载成功")
	return nil
}

// initDatabase 初始化数据库
func initDatabase() error {
	log.Println("正在初始化数据库...")

	// 调用模型初始化
	model.Init()

	log.Println("数据库初始化成功")
	return nil
}

// registerRoutes 注册路由
func registerRoutes() error {
	log.Println("正在注册路由...")

	// 使用统一的路由管理器
	router.Init()

	log.Println("路由注册成功")
	return nil
}

// main 主函数
func main() {
	log.Println("=== 启动Web服务器 ===")

	// 开发模式下开启ORM调试
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
		log.Println("开发模式：ORM调试已开启")
	}

	// 启动服务器
	port, err := beego.AppConfig.String("httpport")
	if err != nil {
		log.Printf("获取端口配置失败，使用默认端口8080")
		port = "8080"
	}
	log.Printf("服务器启动在端口: %s", port)
	beego.Run()
}
