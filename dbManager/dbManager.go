package dbManager

import (
	"fmt"
	"github.com/mrzhangs520/go-tiger/config"
	"github.com/mrzhangs520/go-tiger/dError"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func init() {
	mysqlConfig := config.GetInstance().Section("mysql")
	username := mysqlConfig.Key("username").Value()
	password := mysqlConfig.Key("password").Value()
	host := mysqlConfig.Key("host").Value()
	port := mysqlConfig.Key("port").Value()
	dbname := mysqlConfig.Key("dbname").Value()
	timeout := mysqlConfig.Key("timeout").Value()

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local&timeout=%s", username, password, host, port, dbname, timeout)

	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: new(logger),
	})

	if nil != err {
		panic(dError.NewError("连接数据库出错", err))
	}
	// 控制数据库连接池
	sqlDB, err := db.DB()

	if nil != err {
		panic(dError.NewError("数据库连接池错误", err))
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetInstance() *gorm.DB {
	return db
}
