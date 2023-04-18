package mysql

import (
	"bj38web/web/model"
	"fmt"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var GormDB *gorm.DB

func Init() (err error) {
	// dialector 配置mysql驱动 配置gorm
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")
	charset := viper.GetString("mysql.charset")
	maxOpenConns := viper.GetInt("mysql.max_open_conns")
	maxIdleConns := viper.GetInt("mysql.max_idle_conns")
	//maxConnLifetime := viper.GetInt("mysql.max_conn_lifetime")

	//dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		username, password, host, port, database, charset)
	fmt.Println("dsn", dsn)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		// parseTime 转化时间类型为time.Time loc设置当地时间
		DSN: dsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("数据库连接失败")
		return
	}
	GormDB = gormDB
	// gorm的数据库对象转为 sql的数据库对象 gorm 使用sql包的对象来维护连接池
	sqlDB, _ := GormDB.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//sqlDB.SetConnMaxLifetime(time.Duration(maxConnLifetime) * time.Minute)

	// 连接测试
	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("数据库ping失败")
		return
	}
	fmt.Println("数据库ping成功")

	err = GormDB.AutoMigrate(new(model.User), new(model.House), new(model.Area), new(model.Facility), new(model.HouseImage), new(model.OrderHouse))
	if err != nil {
		fmt.Println("对象迁移数据库表失败")
		return
	}
	fmt.Println("对象迁移数据库表成功")
	return nil
}

// Close 关闭MySQL连接
func Close() {
	db, err := GormDB.DB()
	if err != nil {
		return
	}
	db.Close()
}
