package main

import (
	"log"

	deviceModel "Backend_Lili/internal/device/model"
	priceModel "Backend_Lili/internal/price/model"
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
func loadConfig() error { //函数返回值类型为error接口类型
	log.Println("正在加载配置文件...")

	// 尝试多个可能的配置文件路径
	configPaths := []string{
		"pkg/conf/app.conf",       // 从项目根目录运行
		"../../pkg/conf/app.conf", // 从cmd/api目录运行
		"../pkg/conf/app.conf",    // 从cmd目录运行
		"./pkg/conf/app.conf",     // 当前目录
	}

	var err error
	for _, path := range configPaths {
		err = beego.LoadAppConfig("ini", path)
		if err == nil {
			log.Printf("成功加载配置文件: %s", path)
			break
		}
	}

	if err != nil {
		return err
	}

	// 验证关键配置
	requiredConfigs := []string{"db_host", "db_port", "db_user", "db_password", "db_name"} //创建一个字符串类型的切片 即动态数组
	for _, config := range requiredConfigs {
		if _, err := beego.AppConfig.String(config); err != nil { //如果配置文件中没有这个配置项，则返回错误
			return err
		}
	}

	log.Println("配置文件加载成功")
	return nil
}

// initDatabase 初始化数据库
func initDatabase() error {
	log.Println("正在初始化数据库...")

	// 初始化用户模块数据模型
	model.Init()
	log.Println("用户模块数据模型初始化完成")

	// 初始化设备模块数据模型
	deviceModel.Init()
	log.Println("设备模块数据模型初始化完成")

	// 初始化价格模块数据模型
	priceModel.Init()
	log.Println("价格模块数据模型初始化完成")

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

	// 生产模式下关闭ORM调试
	if beego.BConfig.RunMode == "prod" {
		orm.Debug = false
		log.Println("生产模式：ORM调试已关闭")
	}

	// 启动服务器
	port, err := beego.AppConfig.String("httpport")
	if err != nil {
		log.Printf("获取端口配置失败，使用默认端口8080")
		port = "8080"
	}
	log.Printf("服务器启动在端口: %s", port)

	// 生产环境下监听所有网络接口
	if beego.BConfig.RunMode == "prod" {
		log.Printf("生产模式：服务器将监听所有网络接口 (0.0.0.0:%s)", port)
	}

	beego.Run()
}
