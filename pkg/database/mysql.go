package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// 导入但不直接使用，以下是一些常用的数据库驱动
	_ "github.com/go-sql-driver/mysql"

	"gitee.com/huajinet/go-example/internal/model"
)

func NewGorm() *gorm.DB {
	// DSN 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:changeme@tcp(mysql:3306)/go_dev?charset=utf8mb4&parseTime=True&loc=Local"
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	if err = orm.Migrator().AutoMigrate(
		new(model.Book),
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return orm
}
